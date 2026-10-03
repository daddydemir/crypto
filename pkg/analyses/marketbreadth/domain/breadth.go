package domain

import (
	"math"
	"time"

	madomain "github.com/daddydemir/crypto/pkg/analyses/ma/domain"
	rsidomain "github.com/daddydemir/crypto/pkg/analyses/rsi/domain"
)

type Metric string

const (
	AboveMA7      Metric = "above_ma7"
	AboveMA25     Metric = "above_ma25"
	AboveMA99     Metric = "above_ma99"
	MA7AboveMA25  Metric = "ma7_above_ma25"
	MA25AboveMA99 Metric = "ma25_above_ma99"
	Bullish       Metric = "bullish"
	Bearish       Metric = "bearish"
	RSIAbove50    Metric = "rsi_above_50"
	RSIAbove60    Metric = "rsi_above_60"
	RSIAbove70    Metric = "rsi_above_70"
	RSIBelow30    Metric = "rsi_below_30"
	Positive24H   Metric = "positive_24h"
)

var ValidMetrics = map[Metric]bool{AboveMA7: true, AboveMA25: true, AboveMA99: true, MA7AboveMA25: true, MA25AboveMA99: true, Bullish: true, Bearish: true, RSIAbove50: true, RSIAbove60: true, RSIAbove70: true, RSIBelow30: true, Positive24H: true}

type Candle struct {
	Symbol string
	Time   time.Time
	Close  float64
}
type Coin struct {
	ID, Symbol, Name string
	Rank             int
}
type Ratio struct {
	Count      int     `json:"count"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}
type Participation struct {
	Advancing           int      `json:"advancing"`
	Declining           int      `json:"declining"`
	Unchanged           int      `json:"unchanged"`
	Total               int      `json:"total"`
	AdvanceDeclineRatio *float64 `json:"advanceDeclineRatio"`
	NetAdvance          int      `json:"netAdvance"`
}
type Trend struct {
	AboveMA7      Ratio `json:"aboveMA7"`
	AboveMA25     Ratio `json:"aboveMA25"`
	AboveMA99     Ratio `json:"aboveMA99"`
	MA7AboveMA25  Ratio `json:"ma7AboveMA25"`
	MA25AboveMA99 Ratio `json:"ma25AboveMA99"`
	Bullish       Ratio `json:"bullish"`
	Mixed         Ratio `json:"mixed"`
	Bearish       Ratio `json:"bearish"`
}
type Momentum struct {
	RSIAbove50 Ratio `json:"rsiAbove50"`
	RSIAbove60 Ratio `json:"rsiAbove60"`
	RSIAbove70 Ratio `json:"rsiAbove70"`
	RSIBelow30 Ratio `json:"rsiBelow30"`
}
type ScoreComponents struct {
	Trend         float64 `json:"trend"`
	Momentum      float64 `json:"momentum"`
	Participation float64 `json:"participation"`
}
type Scores struct {
	BreadthScore float64         `json:"breadthScore"`
	Components   ScoreComponents `json:"components"`
}
type CoinMetric struct {
	ID           string   `json:"id"`
	Symbol       string   `json:"symbol"`
	Name         string   `json:"name"`
	CurrentPrice float64  `json:"currentPrice"`
	Change24H    *float64 `json:"change24h"`
	MetricValue  float64  `json:"metricValue"`
}
type Snapshot struct {
	At            time.Time               `json:"at"`
	Participation Participation           `json:"participation"`
	Positive24H   Ratio                   `json:"positive24h"`
	Positive7D    Ratio                   `json:"positive7d"`
	Trend         Trend                   `json:"trend"`
	Momentum      Momentum                `json:"momentum"`
	Scores        Scores                  `json:"scores"`
	Coins         map[Metric][]CoinMetric `json:"-"`
}

func ratio(count, total int) Ratio {
	p := 0.0
	if total > 0 {
		p = float64(count) * 100 / float64(total)
	}
	return Ratio{count, total, round(p)}
}
func round(v float64) float64 { return math.Round(v*100) / 100 }

func Calculate(at time.Time, coins []Coin, series map[string][]Candle, tolerance float64, candlesPerDay int) Snapshot {
	if candlesPerDay < 1 {
		candlesPerDay = 1
	}
	s := Snapshot{At: at, Coins: map[Metric][]CoinMetric{}}
	counts := map[Metric]int{}
	totals := map[Metric]int{}
	mixed := 0
	trendTotal := 0
	for _, coin := range coins {
		all := series[coin.Symbol]
		var values []float64
		for _, c := range all {
			if !c.Time.After(at) {
				values = append(values, c.Close)
			}
		}
		if len(values) == 0 {
			continue
		}
		current := values[len(values)-1]
		var change24, change7 *float64
		if len(values) > candlesPerDay && values[len(values)-1-candlesPerDay] != 0 {
			v := (current - values[len(values)-1-candlesPerDay]) / values[len(values)-1-candlesPerDay] * 100
			change24 = &v
			s.Participation.Total++
			if math.Abs(v) <= tolerance {
				s.Participation.Unchanged++
			} else if v > 0 {
				s.Participation.Advancing++
			} else {
				s.Participation.Declining++
			}
			totals[Positive24H]++
			if v > 0 {
				counts[Positive24H]++
				add(&s, Positive24H, coin, current, change24, v)
			}
		}
		sevenDaySteps := candlesPerDay * 7
		if len(values) > sevenDaySteps && values[len(values)-1-sevenDaySteps] != 0 {
			v := (current - values[len(values)-1-sevenDaySteps]) / values[len(values)-1-sevenDaySteps] * 100
			change7 = &v
			_ = change7
			totals[Metric("positive_7d")]++
			if v > 0 {
				counts[Metric("positive_7d")]++
			}
		}
		if len(values) >= 99 {
			dates := make([]time.Time, len(values))
			pts := madomain.CalculateSeries(dates, values, 7, 25, 99)
			p := pts[len(pts)-1]
			trendTotal++
			checks := map[Metric]bool{AboveMA7: p.Price > p.MA7, AboveMA25: p.Price > p.MA25, AboveMA99: p.Price > p.MA99, MA7AboveMA25: p.MA7 > p.MA25, MA25AboveMA99: p.MA25 > p.MA99, Bullish: p.MA7 > p.MA25 && p.MA25 > p.MA99, Bearish: p.MA7 < p.MA25 && p.MA25 < p.MA99}
			for m, ok := range checks {
				totals[m]++
				if ok {
					counts[m]++
					add(&s, m, coin, current, change24, metricValue(m, p, 0))
				}
			}
			if !checks[Bullish] && !checks[Bearish] {
				mixed++
			}
		}
		if len(values) >= 15 {
			r := rsidomain.Calculate(values)
			for _, m := range []Metric{RSIAbove50, RSIAbove60, RSIAbove70, RSIBelow30} {
				totals[m]++
			}
			checks := map[Metric]bool{RSIAbove50: r > 50, RSIAbove60: r > 60, RSIAbove70: r > 70, RSIBelow30: r < 30}
			for m, ok := range checks {
				if ok {
					counts[m]++
					add(&s, m, coin, current, change24, r)
				}
			}
		}
	}
	if s.Participation.Declining > 0 {
		v := round(float64(s.Participation.Advancing) / float64(s.Participation.Declining))
		s.Participation.AdvanceDeclineRatio = &v
	}
	s.Participation.NetAdvance = s.Participation.Advancing - s.Participation.Declining
	s.Positive24H = ratio(counts[Positive24H], totals[Positive24H])
	s.Positive7D = ratio(counts[Metric("positive_7d")], totals[Metric("positive_7d")])
	s.Trend = Trend{ratio(counts[AboveMA7], totals[AboveMA7]), ratio(counts[AboveMA25], totals[AboveMA25]), ratio(counts[AboveMA99], totals[AboveMA99]), ratio(counts[MA7AboveMA25], totals[MA7AboveMA25]), ratio(counts[MA25AboveMA99], totals[MA25AboveMA99]), ratio(counts[Bullish], trendTotal), ratio(mixed, trendTotal), ratio(counts[Bearish], trendTotal)}
	s.Momentum = Momentum{ratio(counts[RSIAbove50], totals[RSIAbove50]), ratio(counts[RSIAbove60], totals[RSIAbove60]), ratio(counts[RSIAbove70], totals[RSIAbove70]), ratio(counts[RSIBelow30], totals[RSIBelow30])}
	trend := s.Trend.AboveMA25.Percentage*.4 + s.Trend.AboveMA99.Percentage*.3 + s.Trend.Bullish.Percentage*.3
	participation := s.Positive24H.Percentage*.4 + s.Positive7D.Percentage*.6
	momentum := s.Momentum.RSIAbove50.Percentage
	s.Scores = Scores{round(trend*.45 + momentum*.30 + participation*.25), ScoreComponents{round(trend), round(momentum), round(participation)}}
	return s
}

func add(s *Snapshot, m Metric, c Coin, price float64, change *float64, value float64) {
	s.Coins[m] = append(s.Coins[m], CoinMetric{c.ID, c.Symbol, c.Name, price, change, round(value)})
}
func metricValue(m Metric, p madomain.Point, r float64) float64 {
	switch m {
	case AboveMA7:
		return p.MA7
	case AboveMA25:
		return p.MA25
	case AboveMA99:
		return p.MA99
	case MA7AboveMA25, Bullish, Bearish:
		return p.MA7 - p.MA25
	case MA25AboveMA99:
		return p.MA25 - p.MA99
	}
	return r
}
func PercentagePointChange(current, previous float64) float64 { return round(current - previous) }
