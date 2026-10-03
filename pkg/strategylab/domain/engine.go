package domain

import (
	"fmt"
	boll "github.com/daddydemir/crypto/pkg/analyses/bollinger/domain"
	mad "github.com/daddydemir/crypto/pkg/analyses/ma/domain"
	macd "github.com/daddydemir/crypto/pkg/analyses/macd/domain"
	rsi "github.com/daddydemir/crypto/pkg/analyses/rsi/domain"
	"math"
	"time"
)

var fields = map[string]bool{"price": true, "rsi": true, "ma7": true, "ma25": true, "ma99": true, "macd": true, "signal": true, "histogram": true, "upperBand": true, "lowerBand": true, "ma20": true}

func Validate(c Configuration) error {
	for _, g := range []RuleGroup{c.Entry, c.Exit} {
		if g.Logic != LogicAND && g.Logic != LogicOR {
			return fmt.Errorf("logic must be AND or OR")
		}
		if len(g.Conditions) == 0 {
			return fmt.Errorf("rule group must contain conditions")
		}
		for _, r := range g.Conditions {
			if !fields[r.Left] || (r.Operator != "gt" && r.Operator != "lt") {
				return fmt.Errorf("unsupported condition")
			}
			if r.RightField == "" && r.Value == nil {
				return fmt.Errorf("condition requires right field or value")
			}
			if r.RightField != "" && !fields[r.RightField] {
				return fmt.Errorf("unsupported right field")
			}
		}
	}
	return nil
}
func Indicators(points []PricePoint, index int) (map[string]float64, bool) {
	if index < 98 {
		return nil, false
	}
	prices := make([]float64, index+1)
	dates := make([]time.Time, index+1)
	for i := 0; i <= index; i++ {
		prices[i] = points[i].Close
		dates[i] = points[i].Date
	}
	ma := mad.CalculateSeries(dates, prices, 7, 25, 99)
	m := ma[len(ma)-1]
	mc := macd.Calculate(prices, 12, 26, 9)
	bb := boll.CalculateBollinger(prices, dates)
	if len(mc) == 0 || len(bb) == 0 {
		return nil, false
	}
	x := mc[len(mc)-1]
	b := bb[len(bb)-1]
	return map[string]float64{"price": m.Price, "rsi": rsi.Calculate(prices), "ma7": m.MA7, "ma25": m.MA25, "ma99": m.MA99, "macd": x.MACD, "signal": x.Signal, "histogram": x.Histogram, "upperBand": b.UpperBand, "lowerBand": b.LowerBand, "ma20": b.MA20}, true
}
func Evaluate(g RuleGroup, values map[string]float64) (bool, []ConditionResult, bool) {
	results := []ConditionResult{}
	matched := g.Logic == LogicAND
	for _, c := range g.Conditions {
		l, ok := values[c.Left]
		if !ok {
			return false, nil, false
		}
		var r float64
		if c.RightField != "" {
			var found bool
			r, found = values[c.RightField]
			if !found {
				return false, nil, false
			}
		} else if c.Value != nil {
			r = *c.Value
		} else {
			return false, nil, false
		}
		hit := (c.Operator == "gt" && l > r) || (c.Operator == "lt" && l < r)
		results = append(results, ConditionResult{c, l, r, hit})
		if g.Logic == LogicAND {
			matched = matched && hit
		} else {
			matched = matched || hit
		}
	}
	return matched, results, true
}
func MetricsFor(run StrategyRun, trades []StrategyTrade, equity []EquitySnapshot, currentClose float64) Metrics {
	closed := make([]StrategyTrade, 0, len(trades))
	for _, trade := range trades {
		if trade.ExitExecutionDate != nil {
			closed = append(closed, trade)
		}
	}
	trades = closed
	m := Metrics{TotalFees: run.TotalFees, NumberOfTrades: len(trades)}
	m.CurrentEquity = run.Cash + run.Quantity*currentClose
	m.TotalReturn = pct(m.CurrentEquity-run.InitialCapital, run.InitialCapital)
	m.UnrealizedPnL = run.Quantity * (currentClose - run.EntryPrice)
	wins := 0
	profits, losses := 0.0, 0.0
	holding := 0
	for _, t := range trades {
		m.RealizedPnL += t.PnL
		holding += t.HoldingDays
		if t.PnL > 0 {
			wins++
			profits += t.PnL
		} else {
			losses -= t.PnL
		}
		if t.PnL > m.BestTrade || m.NumberOfTrades == 1 {
			m.BestTrade = t.PnL
		}
		if t.PnL < m.WorstTrade || m.NumberOfTrades == 1 {
			m.WorstTrade = t.PnL
		}
	}
	if len(trades) > 0 {
		m.WinRate = float64(wins) * 100 / float64(len(trades))
		m.AverageHoldingDays = float64(holding) / float64(len(trades))
		if wins > 0 {
			m.AverageProfit = profits / float64(wins)
		}
		if len(trades)-wins > 0 {
			m.AverageLoss = -losses / float64(len(trades)-wins)
		}
	}
	if losses > 0 {
		v := profits / losses
		m.ProfitFactor = &v
	}
	for _, e := range equity {
		if e.Drawdown < m.MaxDrawdown {
			m.MaxDrawdown = e.Drawdown
		}
	}
	if len(equity) > 0 {
		m.BenchmarkReturn = pct(equity[len(equity)-1].BenchmarkEquity-run.InitialCapital, run.InitialCapital)
	}
	m.Difference = m.TotalReturn - m.BenchmarkReturn
	return m
}

func NewEquitySnapshot(runID uint, date time.Time, close, equity, benchmarkEquity, initialCapital, previousPeak float64) (EquitySnapshot, float64) {
	peak := math.Max(previousPeak, equity)
	drawdown := 0.0
	if peak > 0 {
		drawdown = (equity - peak) / peak * 100
	}
	return EquitySnapshot{
		RunID: runID, Date: date, Close: close, Equity: equity, BenchmarkEquity: benchmarkEquity,
		StrategyReturn:  pct(equity-initialCapital, initialCapital),
		BenchmarkReturn: pct(benchmarkEquity-initialCapital, initialCapital),
		RunningPeak:     peak, Drawdown: math.Round(drawdown*100) / 100,
	}, peak
}
func pct(v, b float64) float64 {
	if b == 0 {
		return 0
	}
	return math.Round(v/b*10000) / 100
}
func ExecuteBuy(cash, close, feePercent, slippagePercent float64) (quantity, price, fee float64) {
	price = close * (1 + slippagePercent/100)
	fee = cash * feePercent / 100
	quantity = (cash - fee) / price
	return
}
func ExecuteSell(quantity, close, feePercent, slippagePercent float64) (cash, price, fee float64) {
	price = close * (1 - slippagePercent/100)
	gross := quantity * price
	fee = gross * feePercent / 100
	cash = gross - fee
	return
}
func RiskExit(entry, dailyLow, dailyHigh float64, r RiskManagement) bool {
	return RiskExitReason(entry, dailyLow, dailyHigh, r) != ""
}
func RiskExitReason(entry, dailyLow, dailyHigh float64, r RiskManagement) string {
	if entry <= 0 {
		return ""
	}
	// When both levels are touched in the same daily candle their intraday order is
	// unknowable. Returning on the stop first is the conservative assumption.
	if r.StopLossPercent != nil && dailyLow <= entry*(1-*r.StopLossPercent/100) {
		return "STOP_LOSS"
	}
	if r.TakeProfitPercent != nil && dailyHigh >= entry*(1+*r.TakeProfitPercent/100) {
		return "TAKE_PROFIT"
	}
	return ""
}
