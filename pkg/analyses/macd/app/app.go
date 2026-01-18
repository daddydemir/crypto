package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/macd/domain"
	"github.com/daddydemir/crypto/pkg/analyses/macd/infra"
)

type App struct {
	repo *infra.Repository
}

type MACDDTO struct {
	CoinID   string          `json:"coin_id"`
	Name     string          `json:"name"`
	MACDData domain.MACDData `json:"macd_data"`
}

type MACDHistoryDTO struct {
	Date      string  `json:"date"`
	MACD      float64 `json:"macd"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
}

func NewApp(repo *infra.Repository) *App {
	return &App{repo: repo}
}

// CoinMACDHistory returns historical MACD data for a specific coin
func (a *App) CoinMACDHistory(coinID string) ([]MACDHistoryDTO, error) {

	fastPeriod := 12
	slowPeriod := 26
	signalPeriod := 9

	// Get price data from database
	priceData, err := a.repo.GetCoinPricesFromDB(coinID)
	if err != nil {
		return nil, err
	}

	if len(priceData) < slowPeriod+signalPeriod {
		return []MACDHistoryDTO{}, nil
	}

	// Extract prices for calculation
	prices := make([]float64, len(priceData))
	for i, p := range priceData {
		prices[i] = p.Price
	}

	// Calculate MACD
	macdData := domain.Calculate(prices, fastPeriod, slowPeriod, signalPeriod)
	if len(macdData) == 0 {
		return []MACDHistoryDTO{}, nil
	}

	// Convert to DTO with dates
	var results []MACDHistoryDTO
	startIndex := len(priceData) - len(macdData)

	for i, macd := range macdData {
		results = append(results, MACDHistoryDTO{
			Date:      priceData[startIndex+i].Date.Format("2006-01-02"),
			MACD:      macd.MACD,
			Signal:    macd.Signal,
			Histogram: macd.Histogram,
		})
	}

	return results, nil
}
