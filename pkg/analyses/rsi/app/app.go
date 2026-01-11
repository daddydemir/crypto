package app

import "github.com/daddydemir/crypto/pkg/analyses/rsi/infra"

type App struct {
	repo infra.Repository
}

func NewApp(repo infra.Repository) *App {
	return &App{repo: repo}
}

func (a *App) Execute(coinID string, days int) ([]RSIHistoryDTO, error) {
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

		rsi := a.rsiCalc.Calculate(window)
		results = append(results, RSIHistoryDTO{
			Date: priceData[i].Date.Format("2006-01-02"),
			RSI:  rsi,
		})
	}

	return results, nil
}
