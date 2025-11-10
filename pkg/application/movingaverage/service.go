package movingaverage

import (
	"context"
	"github.com/daddydemir/crypto/pkg/domain/movingaverage"
	"github.com/daddydemir/crypto/pkg/infrastructure"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"slices"
	"time"
)

var list = []string{"usd-coin", "paypal-usd", "ripple-usd", "tether", "ethena-usde"}

type Service struct {
	PriceRepo movingaverage.PriceHistoryRepository
	repo      infrastructure.PriceRepository
}

func NewService(priceRepo movingaverage.PriceHistoryRepository, repo infrastructure.PriceRepository) *Service {
	return &Service{
		PriceRepo: priceRepo,
		repo:      repo,
	}
}

func (s *Service) GetMovingAverageSeries(ctx context.Context, coinID string, days int) ([]movingaverage.Point, error) {
	datas, err := s.PriceRepo.GetLastNDaysPricesWithDates(ctx, coinID, days)
	if err != nil {
		return nil, err
	}
	dates, prices := make([]time.Time, 0, len(datas)), make([]float64, 0, len(datas))
	for _, d := range datas {
		dates = append(dates, d.Date)
		prices = append(prices, d.Price)
	}
	return movingaverage.CalculateSeries(dates, prices), nil
}

func (s *Service) GetMovingAverageSignals() ([]movingaverage.Signal, error) {
	coins, err := s.repo.GetTopCoins()
	if err != nil {
		return nil, err
	}
	response := make([]movingaverage.Signal, 0, len(coins))
	for _, c := range coins {
		if slices.Contains(list, c.Id) {
			continue
		}
		prices, err := s.repo.GetHistoricalPrices(c.Id, -101)
		if err != nil {
			continue
		}
		series := movingaverage.CalculateSeries(getTimeAndPrice(prices))
		if len(series) != 3 {
			continue
		}
		signal := series[2]
		plus := signal.MA7 > signal.MA25 && signal.MA25 > signal.MA99
		minus := signal.MA7 < signal.MA25 && signal.MA25 < signal.MA99

		if plus || minus {
			response = append(response, movingaverage.Signal{
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
