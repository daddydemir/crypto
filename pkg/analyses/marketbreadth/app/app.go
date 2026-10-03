package app

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/daddydemir/crypto/pkg/analyses/marketbreadth/domain"
	"github.com/daddydemir/crypto/pkg/analyses/marketbreadth/infra"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/infrastructure"
)

const UniverseLabel = "Top 100 non-stablecoin cryptocurrencies"

var stable = map[string]bool{"USDT": true, "USDC": true, "DAI": true, "FDUSD": true, "USDE": true, "PYUSD": true, "USDD": true, "TUSD": true, "USDP": true, "USD1": true, "USD0": true, "USDG": true, "USDY": true, "RLUSD": true}

type Repository interface {
	Candles([]string, int) (map[string][]domain.Candle, error)
}
type Catalog interface {
	List() ([]infrastructure.Coin, error)
}
type App struct {
	repo    Repository
	catalog Catalog
	cache   cache.Cache
}

func NewApp(repo Repository, catalog Catalog, c cache.Cache) *App {
	return &App{repo: repo, catalog: catalog, cache: c}
}

type Change struct {
	OneDay   float64 `json:"oneDay"`
	SevenDay float64 `json:"sevenDay"`
}
type BTCComparison struct {
	Price           float64 `json:"price"`
	Change24H       float64 `json:"change24h"`
	BreadthChange7D float64 `json:"breadthChange7d"`
	Summary         string  `json:"summary"`
}
type Response struct {
	Universe        string          `json:"universe"`
	Timeframe       string          `json:"timeframe"`
	Snapshot        domain.Snapshot `json:"snapshot"`
	BreadthMomentum Change          `json:"breadthMomentum"`
	BTC             BTCComparison   `json:"btc"`
	Warnings        []string        `json:"warnings"`
	Coverage        Coverage        `json:"coverage"`
}
type Coverage struct {
	UniverseCoins int `json:"universeCoins"`
	MA99Valid     int `json:"ma99Valid"`
	ExcludedCoins int `json:"excludedCoins"`
}
type HistoryPoint struct {
	Date          time.Time `json:"date"`
	BreadthScore  float64   `json:"breadthScore"`
	AboveMA25     float64   `json:"aboveMA25"`
	AboveMA99     float64   `json:"aboveMA99"`
	RSIAbove50    float64   `json:"rsiAbove50"`
	PositiveCoins float64   `json:"positiveCoins"`
	Advancing     int       `json:"advancing"`
	Declining     int       `json:"declining"`
}
type HistoryResponse struct {
	Universe  string         `json:"universe"`
	Timeframe string         `json:"timeframe"`
	Period    string         `json:"period"`
	Points    []HistoryPoint `json:"points"`
	Warnings  []string       `json:"warnings"`
}

func (a *App) universe() ([]domain.Coin, error) {
	items, err := a.catalog.List()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].MarketCapRank < items[j].MarketCapRank })
	out := []domain.Coin{}
	for _, c := range items {
		symbol := strings.ToUpper(c.Symbol)
		if stable[symbol] || c.MarketCapRank <= 0 {
			continue
		}
		out = append(out, domain.Coin{ID: c.ID, Symbol: symbol, Name: c.Name, Rank: c.MarketCapRank})
		if len(out) == 100 {
			break
		}
	}
	return out, nil
}
func sanitize(series map[string][]domain.Candle) int {
	invalid := 0
	max := 36 * time.Hour
	for symbol, rows := range series {
		if len(rows) >= 2 && rows[len(rows)-1].Time.Sub(rows[len(rows)-2].Time) > max {
			delete(series, symbol)
			invalid++
		}
	}
	return invalid
}
func symbols(coins []domain.Coin) []string {
	v := make([]string, len(coins))
	for i, c := range coins {
		v[i] = c.Symbol
	}
	return v
}

func (a *App) Current() (Response, error) {
	coins, err := a.universe()
	if err != nil {
		return Response{}, err
	}
	series, err := a.repo.Candles(symbols(coins), 107)
	if err != nil {
		return Response{}, err
	}
	invalid := sanitize(series)
	at := infra.LatestTime(series)
	s := domain.Calculate(at, coins, series, .0001, 1)
	warnings := warnings(len(coins), s.Trend.AboveMA99.Total, invalid)
	one := historicalSnapshot(at.Add(-24*time.Hour), coins, series)
	seven := historicalSnapshot(at.Add(-7*24*time.Hour), coins, series)
	bc := domain.PercentagePointChange(s.Trend.AboveMA25.Percentage, seven.Trend.AboveMA25.Percentage)
	btc := btcComparison(series["BTC"], s, bc)
	return Response{UniverseLabel, "1d", s, Change{domain.PercentagePointChange(s.Trend.AboveMA25.Percentage, one.Trend.AboveMA25.Percentage), bc}, btc, warnings, Coverage{len(coins), s.Trend.AboveMA99.Total, invalid}}, nil
}
func historicalSnapshot(at time.Time, coins []domain.Coin, series map[string][]domain.Candle) domain.Snapshot {
	return domain.Calculate(at, coins, series, .0001, 1)
}
func warnings(universe, valid, invalid int) []string {
	w := []string{}
	if valid*100 < universe*80 {
		w = append(w, fmt.Sprintf("Limited MA99 coverage: %d of %d universe coins have valid data.", valid, universe))
	}
	if invalid > 0 {
		w = append(w, fmt.Sprintf("%d coins were excluded because daily candle data is unavailable or stale.", invalid))
	}
	return w
}
func btcComparison(rows []domain.Candle, s domain.Snapshot, breadthChange float64) BTCComparison {
	b := BTCComparison{BreadthChange7D: breadthChange}
	if len(rows) > 0 {
		b.Price = rows[len(rows)-1].Close
	}
	if len(rows) > 1 && rows[len(rows)-2].Close != 0 {
		b.Change24H = (b.Price - rows[len(rows)-2].Close) / rows[len(rows)-2].Close * 100
	}
	if b.Change24H > 0 && breadthChange < 0 {
		b.Summary = "BTC is rising while market participation has declined."
	} else if b.Change24H < 0 && breadthChange > 0 {
		b.Summary = "BTC is declining while market participation has improved."
	} else {
		b.Summary = "BTC and market participation are moving in the same direction."
	}
	return b
}

func (a *App) History(period string) (HistoryResponse, error) {
	key := "crypto:market-breadth:history:1d:" + period
	if raw, err := a.cache.Get(key); err == nil {
		if text, ok := raw.(string); ok {
			var cached HistoryResponse
			if json.Unmarshal([]byte(text), &cached) == nil {
				return cached, nil
			}
		}
	}
	coins, err := a.universe()
	if err != nil {
		return HistoryResponse{}, err
	}
	days := map[string]int{"30d": 30, "90d": 90, "1y": 365}[period]
	pointsNeeded := days + 105
	series, err := a.repo.Candles(symbols(coins), pointsNeeded)
	if err != nil {
		return HistoryResponse{}, err
	}
	invalid := sanitize(series)
	at := infra.LatestTime(series)
	step := 24 * time.Hour
	start := at.Add(-time.Duration(days) * 24 * time.Hour)
	points := []HistoryPoint{}
	for cursor := start; !cursor.After(at); cursor = cursor.Add(step) {
		s := historicalSnapshot(cursor, coins, series)
		points = append(points, HistoryPoint{cursor, s.Scores.BreadthScore, s.Trend.AboveMA25.Percentage, s.Trend.AboveMA99.Percentage, s.Momentum.RSIAbove50.Percentage, s.Positive24H.Percentage, s.Participation.Advancing, s.Participation.Declining})
	}
	resp := HistoryResponse{UniverseLabel, "1d", period, points, warnings(len(coins), pointsValid(series), invalid)}
	if encoded, e := json.Marshal(resp); e == nil {
		_ = a.cache.SetWithExpiration(key, encoded, 15*time.Minute)
	}
	return resp, nil
}
func pointsValid(series map[string][]domain.Candle) int {
	n := 0
	for _, v := range series {
		if len(v) >= 99 {
			n++
		}
	}
	return n
}

func (a *App) Coins(metric domain.Metric) ([]domain.CoinMetric, error) {
	current, err := a.Current()
	if err != nil {
		return nil, err
	}
	return current.Snapshot.Coins[metric], nil
}
