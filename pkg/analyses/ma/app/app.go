package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/ma/domain"
	"github.com/daddydemir/crypto/pkg/analyses/ma/infra"
	"github.com/daddydemir/crypto/pkg/infrastructure"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"slices"
	"time"
)

var list = []string{"usd-coin", "paypal-usd", "ripple-usd", "tether", "ethena-usde"}

type App struct {
	PriceRepo *infra.Repository
	repo      infrastructure.PriceRepository
}

func NewApp(priceRepo *infra.Repository, repo infrastructure.PriceRepository) *App {
	return &App{
		PriceRepo: priceRepo,
		repo:      repo,
	}
}

func (a *App) GetMovingAverageSeries(coinID string, days int) ([]domain.Point, error) {
	datas, err := a.PriceRepo.GetLastNDaysPricesWithDates(coinID, days)
	if err != nil {
		return nil, err
	}
	dates, prices := make([]time.Time, 0, len(datas)), make([]float64, 0, len(datas))
	for _, d := range datas {
		dates = append(dates, d.Date)
		prices = append(prices, d.Price)
	}
	return domain.CalculateSeries(dates, prices), nil
}

func (a *App) GetMovingAverageSignals() ([]domain.Signal, error) {
	coins, err := a.repo.GetTopCoins()
	if err != nil {
		return nil, err
	}
	response := make([]domain.Signal, 0, len(coins))
	for _, c := range coins {
		if slices.Contains(list, c.Id) {
			continue
		}
		prices, err := a.repo.GetHistoricalPrices(c.Id, -101)
		if err != nil {
			continue
		}
		series := domain.CalculateSeries(getTimeAndPrice(prices))
		if len(series) != 3 {
			continue
		}
		signal := series[2]
		plus := signal.MA7 > signal.MA25 && signal.MA25 > signal.MA99
		minus := signal.MA7 < signal.MA25 && signal.MA25 < signal.MA99

		if plus || minus {
			response = append(response, domain.Signal{
				Id:     c.Id,
				Name:   c.Name,
				Symbol: c.Symbol,
				Price:  float64(c.PriceUsd),
				Points: series,
			})
		}
	}
	return response, nil
}

func getTimeAndPrice(list []coincap.History) ([]time.Time, []float64) {
	dates, prices := make([]time.Time, 0, len(list)), make([]float64, 0, len(list))
	for _, h := range list {
		dates = append(dates, h.Date)
		prices = append(prices, float64(h.PriceUsd))
	}
	return dates, prices
}
