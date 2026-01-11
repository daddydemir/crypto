package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/bollinger/domain"
	"github.com/daddydemir/crypto/pkg/infrastructure"

	"time"
)

type App struct {
	PriceRepo domain.Repository
	repo      infrastructure.PriceRepository
}

func NewApp(priceRepo domain.Repository, repo infrastructure.PriceRepository) *App {
	return &App{
		PriceRepo: priceRepo,
		repo:      repo,
	}
}

func (a *App) GetBollingerSeries(coinID string, days int) ([]domain.Point, error) {
	datas, err := a.PriceRepo.GetLastNDaysPricesWithDates(coinID, days)
	if err != nil {
		return nil, err
	}
	prices, dates := convert(datas)
	return domain.CalculateBollinger(prices, dates), nil
}

func (a *App) GetBollingerBandSignal() ([]domain.Signal, error) {
	coins, err := a.repo.GetTopCoins()
	if err != nil {
		return nil, err
	}

	response := make([]domain.Signal, 0, len(coins))

	for _, coin := range coins {
		prices, err := a.PriceRepo.GetLastNDaysPricesWithDates(coin.Id, -20)
		if err != nil || len(prices) != 20 {
			continue
		}
		floats, times := convert(prices)
		points := domain.CalculateBollinger(floats, times)

		if points[0].LowerBand > prices[19].Price || points[0].UpperBand < prices[19].Price {
			response = append(response, domain.Signal{
				Id:     coin.Id,
				Name:   coin.Name,
				Symbol: coin.Symbol,
				Price:  prices[19].Price,
				Point:  points[0],
			})
		}
	}

	return response, nil
}

func convert(list []domain.PriceData) ([]float64, []time.Time) {
	dates, prices := make([]time.Time, 0, len(list)), make([]float64, 0, len(list))
	for _, d := range list {
		dates = append(dates, d.Date)
		prices = append(prices, d.Price)
	}

	return prices, dates
}
