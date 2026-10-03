package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/pkg/strategylab/domain"
	"github.com/daddydemir/crypto/pkg/strategylab/infra"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"math"
	"strconv"
	"time"
)

const backtestQueue = "strategy-lab:backtests"

type App struct {
	r     *infra.Repository
	redis *redis.Client
}

func New(r *infra.Repository, client *redis.Client) *App      { return &App{r: r, redis: client} }
func (a *App) List(u string) ([]domain.Strategy, error)       { return a.r.ListStrategies(u) }
func (a *App) Get(u string, id uint) (domain.Strategy, error) { return a.r.Strategy(u, id) }
func (a *App) Create(u string, s domain.Strategy) (domain.Strategy, error) {
	if err := domain.Validate(s.Configuration); err != nil {
		return s, err
	}
	s.ID = 0
	s.Username = u
	s.Version = 1
	s.RunCount = 0
	return s, a.r.SaveStrategy(&s)
}
func (a *App) Update(u string, id uint, s domain.Strategy) (domain.Strategy, error) {
	old, e := a.r.Strategy(u, id)
	if e != nil {
		return s, e
	}
	if e = domain.Validate(s.Configuration); e != nil {
		return s, e
	}
	old.Name = s.Name
	old.Description = s.Description
	old.Configuration = s.Configuration
	old.Version++
	return old, a.r.SaveStrategy(&old)
}
func (a *App) Delete(u string, id uint) error { return a.r.DeleteStrategy(u, id) }
func (a *App) Clone(u string, id uint) (domain.Strategy, error) {
	s, e := a.r.Strategy(u, id)
	if e != nil {
		return s, e
	}
	s.ID = 0
	s.Name += " Copy"
	s.Version = 1
	s.RunCount = 0
	s.CreatedAt = time.Time{}
	s.UpdatedAt = time.Time{}
	return s, a.r.SaveStrategy(&s)
}

type RunRequest struct {
	CoinSymbol      string     `json:"coinSymbol"`
	StartDate       time.Time  `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	InitialCapital  float64    `json:"initialCapital"`
	FeePercent      float64    `json:"feePercent"`
	SlippagePercent float64    `json:"slippagePercent"`
}

func (a *App) Start(u string, id uint, typ string, q RunRequest) (domain.StrategyRun, error) {
	s, e := a.r.Strategy(u, id)
	if e != nil {
		return domain.StrategyRun{}, e
	}
	if q.InitialCapital <= 0 {
		return domain.StrategyRun{}, fmt.Errorf("initial capital must be positive")
	}
	run := domain.StrategyRun{Username: u, StrategyID: id, StrategySnapshot: domain.StrategySnapshot{id, s.Name, s.Version, s.Configuration}, CoinSymbol: q.CoinSymbol, Type: typ, Status: domain.StatusRunning, StartDate: q.StartDate, EndDate: q.EndDate, InitialCapital: q.InitialCapital, FeePercent: q.FeePercent, SlippagePercent: q.SlippagePercent, Cash: q.InitialCapital}
	if typ == domain.RunBacktest {
		run.Status = domain.StatusQueued
	}
	if e = a.r.SaveRun(&run); e != nil {
		return run, e
	}
	if e = a.r.IncrementRunCount(id); e != nil {
		return run, e
	}
	if typ == domain.RunBacktest {
		e = a.redis.RPush(context.Background(), backtestQueue, run.ID).Err()
	} else {
		e = a.process(&run, false)
	}
	return run, e
}

func (a *App) StartBacktestWorker() {
	go func() {
		for {
			item, err := a.redis.BLPop(context.Background(), 0, backtestQueue).Result()
			if err != nil {
				slog.Error("strategy backtest queue", "error", err)
				time.Sleep(time.Second)
				continue
			}
			id, err := strconv.ParseUint(item[1], 10, 64)
			if err != nil {
				continue
			}
			run, err := a.r.RunByID(uint(id))
			if err != nil || run.Status != domain.StatusQueued {
				continue
			}
			run.Status = domain.StatusRunning
			_ = a.r.SaveRun(&run)
			err = a.process(&run, true)
			if err != nil {
				run.Status = domain.StatusStopped
				slog.Error("strategy backtest", "run", run.ID, "error", err)
			} else {
				run.Status = domain.StatusCompleted
			}
			_ = a.r.SaveRun(&run)
			a.notifyCompleted(run)
		}
	}()
}

func (a *App) notifyCompleted(run domain.StrategyRun) {
	payload, _ := json.Marshal(map[string]any{"username": run.Username, "tip": "STRATEGY", "coin": run.CoinSymbol + " · " + run.StrategySnapshot.Name, "gorsel": "", "olusturma_zamani": time.Now().Unix()})
	a.redis.ZAdd(context.Background(), "kripto:grafikler", redis.Z{Score: float64(time.Now().Add(24 * time.Hour).Unix()), Member: payload}).Err()
}
func (a *App) process(run *domain.StrategyRun, backtest bool) error {
	points, e := a.r.Prices(run.CoinSymbol)
	if e != nil {
		return e
	}
	benchmarkQty := run.BenchmarkQuantity
	runningPeak := 0.0
	if !backtest {
		runningPeak, e = a.r.PeakEquity(run.ID)
		if e != nil {
			return e
		}
	}
	snapshots := make([]domain.EquitySnapshot, 0, len(points))
	for i, p := range points {
		if p.Date.Before(run.StartDate) || run.EndDate != nil && p.Date.After(*run.EndDate) {
			continue
		}
		if run.LastEvaluatedDate != nil && !p.Date.After(*run.LastEvaluatedDate) {
			continue
		}
		if !backtest && !a.r.Mark(run.ID, p.Date) {
			continue
		}
		if benchmarkQty == 0 {
			benchmarkQty = run.InitialCapital * (1 - run.FeePercent/100) / (p.Close * (1 + run.SlippagePercent/100))
			run.BenchmarkQuantity = benchmarkQty
		}
		if run.PendingAction != "" {
			if run.PendingAction == "BUY" && run.Quantity == 0 {
				quantity, price, fee := domain.ExecuteBuy(run.Cash, p.Close, run.FeePercent, run.SlippagePercent)
				run.Quantity = quantity
				run.EntryPrice = price
				run.Cash = 0
				d := p.Date
				run.EntryExecutionDate = &d
				run.TotalFees += fee
				a.r.CreateTrade(&domain.StrategyTrade{RunID: run.ID, EntrySignalDate: *run.PendingSignalDate, EntryExecutionDate: d, EntryPrice: price, EntryFee: fee, EntrySnapshot: run.PendingSnapshot, Quantity: run.Quantity})
			} else if run.PendingAction == "SELL" && run.Quantity > 0 {
				cash, price, fee := domain.ExecuteSell(run.Quantity, p.Close, run.FeePercent, run.SlippagePercent)
				open, _ := a.r.OpenTrade(run.ID)
				run.Cash = cash
				run.TotalFees += fee
				d := p.Date
				open.ExitExecutionDate = &d
				open.ExitSignalDate = run.PendingSignalDate
				open.ExitPrice = price
				open.ExitFee = fee
				open.ExitSnapshot = run.PendingSnapshot
				open.ExitReason = run.PendingExitReason
				open.PnL = run.Cash - (open.Quantity*open.EntryPrice + open.EntryFee)
				open.PnLPercent = open.PnL / (open.Quantity*open.EntryPrice + open.EntryFee) * 100
				open.HoldingDays = int(d.Sub(open.EntryExecutionDate).Hours() / 24)
				a.r.SaveTrade(&open)
				run.Quantity = 0
				run.EntryPrice = 0
				run.EntryExecutionDate = nil
			}
			run.PendingAction = ""
			run.PendingExitReason = ""
			run.PendingSignalDate = nil
		}
		values, ok := domain.Indicators(points, i)
		if ok {
			group := run.StrategySnapshot.Configuration.Entry
			action := "BUY"
			if run.Quantity > 0 {
				group = run.StrategySnapshot.Configuration.Exit
				action = "SELL"
			}
			matched, res, _ := domain.Evaluate(group, values)
			exitReason := "STRATEGY_RULE"
			if run.Quantity > 0 {
				risk := run.StrategySnapshot.Configuration.Risk
				if reason := domain.RiskExitReason(run.EntryPrice, p.Low, p.High, risk); reason != "" {
					matched = true
					exitReason = reason
				}
			}
			if matched {
				d := p.Date
				run.PendingAction = action
				if action == "SELL" {
					run.PendingExitReason = exitReason
				}
				run.PendingSignalDate = &d
				run.PendingSnapshot = domain.SignalSnapshot{Date: d, Values: values, Conditions: res}
			}
		}
		eq := run.Cash + run.Quantity*p.Close
		var snapshot domain.EquitySnapshot
		snapshot, runningPeak = domain.NewEquitySnapshot(run.ID, p.Date, p.Close, eq, benchmarkQty*p.Close, run.InitialCapital, runningPeak)
		if backtest {
			snapshots = append(snapshots, snapshot)
		} else if err := a.r.Snapshot(snapshot); err != nil {
			return err
		}
		d := p.Date
		run.LastEvaluatedDate = &d
	}
	if backtest {
		if err := a.r.Snapshots(snapshots); err != nil {
			return err
		}
	}
	if !backtest && run.EndDate != nil && run.LastEvaluatedDate != nil && !run.LastEvaluatedDate.Before(*run.EndDate) {
		run.Status = domain.StatusCompleted
	}
	if e := a.refreshMetrics(run); e != nil {
		return e
	}
	return a.r.SaveRun(run)
}
func (a *App) Runs(u string) ([]domain.StrategyRun, error) {
	return a.r.Runs(u)
}
func (a *App) Run(u string, id uint) (domain.StrategyRun, error) {
	return a.r.Run(u, id)
}
func (a *App) DeleteRun(u string, id uint) error { return a.r.DeleteRun(u, id) }
func (a *App) refreshMetrics(v *domain.StrategyRun) error {
	t, err := a.r.Trades(v.ID)
	if err != nil {
		return err
	}
	q, err := a.r.Equity(v.ID)
	if err != nil {
		return err
	}
	close := 0.0
	if len(q) > 0 {
		close = q[len(q)-1].Close
	}
	v.Metrics = domain.MetricsFor(*v, t, q, close)
	now := time.Now()
	v.MetricsCalculatedAt = &now
	return nil
}
func (a *App) Trades(u string, id uint) ([]domain.StrategyTrade, error) {
	if _, e := a.r.Run(u, id); e != nil {
		return nil, e
	}
	return a.r.Trades(id)
}
func (a *App) Equity(u string, id uint) ([]domain.EquitySnapshot, error) {
	if _, e := a.r.Run(u, id); e != nil {
		return nil, e
	}
	return a.r.Equity(id)
}
func (a *App) Compare(u string, ids []uint) ([]domain.StrategyRun, error) {
	out := make([]domain.StrategyRun, 0, len(ids))
	for _, id := range ids {
		v, e := a.Run(u, id)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, nil
}
func (a *App) Status(u string, id uint, status string) error {
	v, e := a.r.Run(u, id)
	if e != nil {
		return e
	}
	v.Status = status
	if status == domain.StatusRunning {
		return a.process(&v, false)
	}
	return a.r.SaveRun(&v)
}
func (a *App) EvaluateRunning() {
	runs, _ := a.r.Running()
	for i := range runs {
		_ = a.process(&runs[i], false)
	}
}

var _ = math.Abs
