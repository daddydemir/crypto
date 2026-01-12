package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/ema/domain"
	"time"
)

type App struct {
	repo domain.Repository
}

func NewApp(repo domain.Repository) *App {
	return &App{
		repo: repo,
	}
}

func (a *App) GetMovingAverageSeries(coinId string, days int) ([]domain.Point, error) {
	priceData, err := a.repo.GetLastNDaysPricesWithDates(coinId, days)
	if err != nil {
		return nil, err
	}
	dates, prices := make([]time.Time, 0, len(priceData)), make([]float64, 0, len(priceData))
	for _, d := range priceData {
		dates = append(dates, d.Date)
		prices = append(prices, d.Price)
	}
	return domain.CalculateSeries(dates, prices), nil
}
