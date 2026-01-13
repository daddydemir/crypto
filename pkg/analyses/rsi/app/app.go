package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/rsi/domain"
	"github.com/daddydemir/crypto/pkg/analyses/rsi/infra"
	"strings"
)

type App struct {
	repo *infra.Repository
}

func NewApp(repo *infra.Repository) *App {
	return &App{repo: repo}
}

func (a *App) TopCoinsRSI() ([]RSIDTO, error) {
	coins, err := a.repo.GetTopCoinIDs()
	if err != nil {
		return nil, err
	}
	coinIDs := make([]string, 0, len(coins))
	for _, c := range coins {
		coinIDs = append(coinIDs, strings.ToLower(c.Symbol))
	}
	prices, err := a.repo.GetLastNDaysPrices(coinIDs, 14)
	if err != nil {
		return nil, err
	}

	var dtos []RSIDTO
	for _, c := range coins {
		id := strings.ToLower(c.Symbol)
		rsi := domain.Calculate(prices[id])
		dtos = append(dtos, RSIDTO{
			CoinID: id,
			Name:   c.Name,
			RSI:    rsi,
		})
	}

	return dtos, nil
}

func (a *App) CoinRSIHistory(coinID string, days int) ([]RSIHistoryDTO, error) {
	priceData, err := a.repo.GetHistoricalPrices(coinID, days)
	if err != nil {
		return nil, err
	}

	if len(priceData) < 15 {
		return nil, nil
	}

	var results []RSIHistoryDTO
	for i := 14; i < len(priceData); i++ {
		var window []float64
		for j := i - 14; j <= i; j++ {
			window = append(window, priceData[j].Price)
		}

		rsi := domain.Calculate(window)
		results = append(results, RSIHistoryDTO{
			Date: priceData[i].Date.Format("2006-01-02"),
			RSI:  rsi,
		})
	}

	return results, nil
}
