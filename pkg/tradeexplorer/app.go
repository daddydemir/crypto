package tradeexplorer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	strategy "github.com/daddydemir/crypto/pkg/strategylab/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	db    *gorm.DB
	redis *redis.Client
}

func New(db *gorm.DB, redis *redis.Client) *App { return &App{db: db, redis: redis} }

type Filters struct {
	Result, Coin, RunType string
	StrategyID, RunID     uint
	From, To              *time.Time
	Sort                  string
	Desc                  bool
	Page, PageSize        int
}
type TradeRow struct {
	ID                 uint       `json:"id"`
	RunID              uint       `json:"runId"`
	StrategyID         uint       `json:"strategyId"`
	CoinSymbol         string     `json:"coinSymbol"`
	StrategyName       string     `json:"strategyName"`
	StrategyVersion    int        `json:"strategyVersion"`
	RunType            string     `json:"runType"`
	EntrySignalDate    time.Time  `json:"entrySignalDate"`
	EntryExecutionDate time.Time  `json:"entryExecutionDate"`
	ExitSignalDate     *time.Time `json:"exitSignalDate"`
	ExitExecutionDate  *time.Time `json:"exitExecutionDate"`
	EntryPrice         float64    `json:"entryPrice"`
	ExitPrice          float64    `json:"exitPrice"`
	PnL                float64    `json:"pnl"`
	PnLPercent         float64    `json:"pnlPercent"`
	HoldingDays        int        `json:"holdingDays"`
	TotalFees          float64    `json:"totalFees"`
	ExitReason         string     `json:"exitReason"`
	Note               string     `json:"note"`
}
type Page struct {
	Items    []TradeRow `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
}
type Point struct {
	Date          time.Time `json:"date"`
	Close         float64   `json:"close"`
	ReturnPercent float64   `json:"returnPercent"`
	Phase         string    `json:"phase"`
}
type PostExit struct {
	Days          int      `json:"days"`
	ReturnPercent *float64 `json:"returnPercent"`
}
type Excursion struct {
	BestPercent  float64    `json:"bestPercent"`
	WorstPercent float64    `json:"worstPercent"`
	BestDate     *time.Time `json:"bestDate"`
	WorstDate    *time.Time `json:"worstDate"`
}
type Analysis struct {
	Chart                   []Point    `json:"chart"`
	Journey                 []Point    `json:"journey"`
	Excursion               Excursion  `json:"excursion"`
	PostExit                []PostExit `json:"postExit"`
	BestCloseAfter14        *float64   `json:"bestCloseAfter14"`
	BestCloseAfter14Percent *float64   `json:"bestCloseAfter14Percent"`
	CurrentClose            *float64   `json:"currentClose"`
	UnrealizedPnLPercent    *float64   `json:"unrealizedPnlPercent"`
	DaysOpen                int        `json:"daysOpen"`
	CalculatedAsOf          *time.Time `json:"calculatedAsOf"`
}
type Detail struct {
	Trade            TradeRow                `json:"trade"`
	EntrySnapshot    strategy.SignalSnapshot `json:"entrySnapshot"`
	ExitSnapshot     strategy.SignalSnapshot `json:"exitSnapshot"`
	Analysis         Analysis                `json:"analysis"`
	DailyCloseNotice string                  `json:"dailyCloseNotice"`
}
type Statistics struct {
	Total              int64   `json:"total"`
	Winners            int64   `json:"winners"`
	Losers             int64   `json:"losers"`
	Open               int64   `json:"open"`
	WinRate            float64 `json:"winRate"`
	AverageReturn      float64 `json:"averageReturn"`
	AverageWinner      float64 `json:"averageWinner"`
	AverageLoser       float64 `json:"averageLoser"`
	AverageHoldingDays float64 `json:"averageHoldingDays"`
	AverageFees        float64 `json:"averageFees"`
}
type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
type Options struct {
	Coins      []Option `json:"coins"`
	Strategies []Option `json:"strategies"`
	Runs       []Option `json:"runs"`
}

func base(db *gorm.DB, username string) *gorm.DB {
	return db.Table("strategy_trades t").Joins("join strategy_runs r on r.id=t.run_id").Where("r.username=?", username)
}
func filtered(q *gorm.DB, f Filters) *gorm.DB {
	if f.Result == "winner" {
		q = q.Where("t.exit_execution_date is not null and t.pn_l > 0")
	}
	if f.Result == "loser" {
		q = q.Where("t.exit_execution_date is not null and t.pn_l <= 0")
	}
	if f.Result == "open" {
		q = q.Where("t.exit_execution_date is null")
	}
	if f.Result == "closed" {
		q = q.Where("t.exit_execution_date is not null")
	}
	if f.Coin != "" {
		q = q.Where("upper(r.coin_symbol)=upper(?)", f.Coin)
	}
	if f.StrategyID > 0 {
		q = q.Where("r.strategy_id=?", f.StrategyID)
	}
	if f.RunID > 0 {
		q = q.Where("r.id=?", f.RunID)
	}
	if f.RunType != "" {
		q = q.Where("r.type=?", strings.ToUpper(f.RunType))
	}
	if f.From != nil {
		q = q.Where("t.entry_execution_date>=?", *f.From)
	}
	if f.To != nil {
		q = q.Where("t.entry_execution_date<=?", *f.To)
	}
	return q
}
func (a *App) List(u string, f Filters) (Page, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 25
	}
	q := filtered(base(a.db, u), f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return Page{}, err
	}
	orders := map[string]string{"entry": "t.entry_execution_date", "exit": "t.exit_execution_date", "pnl": "t.pn_l", "pnlPercent": "t.pn_l_percent", "holding": "t.holding_days"}
	col := orders[f.Sort]
	if col == "" {
		col = "t.entry_execution_date"
	}
	dir := " asc"
	if f.Desc || f.Sort == "" {
		dir = " desc"
	}
	items := make([]TradeRow, 0)
	err := q.Select(`t.id,t.run_id,r.strategy_id,r.coin_symbol,r.strategy_snapshot->>'name' as strategy_name,coalesce((r.strategy_snapshot->>'version')::int,0) strategy_version,r.type run_type,t.entry_signal_date,t.entry_execution_date,t.exit_signal_date,t.exit_execution_date,t.entry_price,t.exit_price,t.pn_l,t.pn_l_percent,t.holding_days,(t.entry_fee+t.exit_fee) total_fees,t.exit_reason,t.note`).Order(col + dir).Limit(f.PageSize).Offset((f.Page - 1) * f.PageSize).Scan(&items).Error
	return Page{Items: items, Total: total, Page: f.Page, PageSize: f.PageSize}, err
}
func (a *App) row(u string, id uint) (TradeRow, strategy.StrategyTrade, error) {
	var row TradeRow
	var trade strategy.StrategyTrade
	err := base(a.db, u).Select(`t.id,t.run_id,r.strategy_id,r.coin_symbol,r.strategy_snapshot->>'name' as strategy_name,coalesce((r.strategy_snapshot->>'version')::int,0) strategy_version,r.type run_type,t.entry_signal_date,t.entry_execution_date,t.exit_signal_date,t.exit_execution_date,t.entry_price,t.exit_price,t.pn_l,t.pn_l_percent,t.holding_days,(t.entry_fee+t.exit_fee) total_fees,t.exit_reason,t.note`).Where("t.id=?", id).Scan(&row).Error
	if err == nil && row.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err == nil {
		err = a.db.First(&trade, id).Error
	}
	return row, trade, err
}

type candle struct {
	Date  time.Time
	Close float64
}

func (a *App) Detail(u string, id uint) (Detail, error) {
	row, trade, err := a.row(u, id)
	if err != nil {
		return Detail{}, err
	}
	var latest *time.Time
	a.db.Raw("select max(candle_date) from yahoo_candles where upper(symbol)=upper(?) and close_price is not null", row.CoinSymbol).Scan(&latest)
	key := "trade-explorer:" + strconv.FormatUint(uint64(id), 10) + ":"
	if latest != nil {
		key += latest.Format("20060102")
	}
	var analysis Analysis
	if raw, e := a.redis.Get(context.Background(), key).Bytes(); e == nil && json.Unmarshal(raw, &analysis) == nil {
		return Detail{row, trade.EntrySnapshot, trade.ExitSnapshot, analysis, "Calculated using daily closing prices"}, nil
	}
	analysis, err = a.analyze(row, latest)
	if err != nil {
		return Detail{}, err
	}
	raw, _ := json.Marshal(analysis)
	a.redis.Set(context.Background(), key, raw, 30*time.Minute)
	return Detail{row, trade.EntrySnapshot, trade.ExitSnapshot, analysis, "Calculated using daily closing prices"}, nil
}
func pct(current, entry float64) float64 {
	if entry == 0 {
		return 0
	}
	return math.Round((current-entry)/entry*10000) / 100
}
func (a *App) analyze(t TradeRow, latest *time.Time) (Analysis, error) {
	end := time.Now()
	if t.ExitExecutionDate != nil {
		end = t.ExitExecutionDate.AddDate(0, 0, 45)
	}
	start := t.EntryExecutionDate.AddDate(0, 0, -45)
	var candles []candle
	err := a.db.Raw("select candle_date date,close_price close from yahoo_candles where upper(symbol)=upper(?) and candle_date between ? and ? and close_price is not null order by candle_date", t.CoinSymbol, start, end).Scan(&candles).Error
	if err != nil {
		return Analysis{}, err
	}
	return calculateAnalysis(t, candles, latest, time.Now()), nil
}
func calculateAnalysis(t TradeRow, candles []candle, latest *time.Time, now time.Time) Analysis {
	out := Analysis{CalculatedAsOf: latest, Chart: []Point{}, Journey: []Point{}, PostExit: []PostExit{}}
	best, worst := -math.MaxFloat64, math.MaxFloat64
	var exitIndex = -1
	for i, c := range candles {
		phase := "context"
		if !c.Date.Before(t.EntryExecutionDate) && (t.ExitExecutionDate == nil || !c.Date.After(*t.ExitExecutionDate)) {
			phase = "trade"
		}
		p := Point{c.Date, c.Close, pct(c.Close, t.EntryPrice), phase}
		out.Chart = append(out.Chart, p)
		if phase == "trade" {
			out.Journey = append(out.Journey, p)
			if p.ReturnPercent > best {
				best = p.ReturnPercent
				d := c.Date
				out.Excursion.BestDate = &d
			}
			if p.ReturnPercent < worst {
				worst = p.ReturnPercent
				d := c.Date
				out.Excursion.WorstDate = &d
			}
		}
		if t.ExitExecutionDate != nil && c.Date.Equal(*t.ExitExecutionDate) {
			exitIndex = i
		}
	}
	if best != -math.MaxFloat64 {
		out.Excursion.BestPercent = best
		out.Excursion.WorstPercent = worst
	}
	if t.ExitExecutionDate == nil && len(candles) > 0 {
		c := candles[len(candles)-1].Close
		out.CurrentClose = &c
		v := pct(c, t.EntryPrice)
		out.UnrealizedPnLPercent = &v
		out.DaysOpen = int(now.Sub(t.EntryExecutionDate).Hours() / 24)
	}
	if exitIndex >= 0 {
		bestClose := 0.0
		for _, d := range []int{1, 3, 7, 14, 30} {
			var v *float64
			if exitIndex+d < len(candles) {
				x := pct(candles[exitIndex+d].Close, t.ExitPrice)
				v = &x
			}
			out.PostExit = append(out.PostExit, PostExit{d, v})
		}
		for i := exitIndex + 1; i < len(candles) && i <= exitIndex+14; i++ {
			if candles[i].Close > bestClose {
				bestClose = candles[i].Close
			}
		}
		if bestClose > 0 {
			out.BestCloseAfter14 = &bestClose
			v := pct(bestClose, t.ExitPrice)
			out.BestCloseAfter14Percent = &v
		}
	}
	return out
}
func (a *App) Statistics(u string, f Filters) (Statistics, error) {
	var s Statistics
	err := filtered(base(a.db, u), f).Select(`count(*) total,count(*) filter(where t.pn_l>0 and t.exit_execution_date is not null) winners,count(*) filter(where t.pn_l<=0 and t.exit_execution_date is not null) losers,count(*) filter(where t.exit_execution_date is null) open,coalesce(avg(t.pn_l_percent) filter(where t.exit_execution_date is not null),0) average_return,coalesce(avg(t.pn_l_percent) filter(where t.pn_l>0 and t.exit_execution_date is not null),0) average_winner,coalesce(avg(t.pn_l_percent) filter(where t.pn_l<=0 and t.exit_execution_date is not null),0) average_loser,coalesce(avg(t.holding_days) filter(where t.exit_execution_date is not null),0) average_holding_days,coalesce(avg(t.entry_fee+t.exit_fee),0) average_fees`).Scan(&s).Error
	if s.Winners+s.Losers > 0 {
		s.WinRate = float64(s.Winners) * 100 / float64(s.Winners+s.Losers)
	}
	return s, err
}
func (a *App) Options(u string) (Options, error) {
	o := Options{Coins: []Option{}, Strategies: []Option{}, Runs: []Option{}}
	err := a.db.Raw(`select distinct upper(coin_symbol) value,upper(coin_symbol) label from strategy_runs where username=? order by 1`, u).Scan(&o.Coins).Error
	if err != nil {
		return o, err
	}
	err = a.db.Raw(`select distinct strategy_id::text value,strategy_snapshot->>'name' label from strategy_runs where username=? order by 2`, u).Scan(&o.Strategies).Error
	if err != nil {
		return o, err
	}
	err = a.db.Raw(`select id::text value,concat(coin_symbol,' · ',strategy_snapshot->>'name',' #',id) label from strategy_runs where username=? order by id desc limit 200`, u).Scan(&o.Runs).Error
	return o, err
}
func (a *App) UpdateNote(u string, id uint, note string) error {
	if len(note) > 2000 {
		return fmt.Errorf("note is too long")
	}
	res := a.db.Model(&strategy.StrategyTrade{}).Where("id=? and run_id in (select id from strategy_runs where username=?)", id, u).Update("note", note)
	if res.Error == nil && res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}
