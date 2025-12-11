package bollinger

import (
	"context"
	"github.com/daddydemir/crypto/pkg/domain/bollinger"
	"github.com/daddydemir/crypto/pkg/infrastructure"

	"time"
)

type Service struct {
	PriceRepo bollinger.PriceHistoryRepository
	repo      infrastructure.PriceRepository
}

func NewService(priceRepo bollinger.PriceHistoryRepository, repo infrastructure.PriceRepository) *Service {
	return &Service{
		PriceRepo: priceRepo,
		repo:      repo,
	}
}

func (s *Service) GetBollingerSeries(_ context.Context, coinID string, days int) ([]bollinger.Point, error) {
	datas, err := s.PriceRepo.GetLastNDaysPricesWithDates(coinID, days)
	if err != nil {
		return nil, err
	}
	prices, dates := convert(datas)
	return bollinger.CalculateBollinger(prices, dates), nil
}

func (s *Service) GetBollingerBandSignal() ([]bollinger.Signal, error) {
	coins, err := s.repo.GetTopCoins()
	if err != nil {
		return nil, err
	}

	response := make([]bollinger.Signal, 0, len(coins))

	for _, coin := range coins {
		prices, err := s.PriceRepo.GetLastNDaysPricesWithDates(coin.Id, -20)
		if err != nil || len(prices) != 20 {
			continue
		}
		floats, times := convert(prices)
		points := bollinger.CalculateBollinger(floats, times)

		if points[0].LowerBand > prices[19].Price || points[0].UpperBand < prices[19].Price {
			response = append(response, bollinger.Signal{
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

func convert(list []bollinger.PriceData) ([]float64, []time.Time) {
	dates, prices := make([]time.Time, 0, len(list)), make([]float64, 0, len(list))
	for _, d := range list {
		dates = append(dates, d.Date)
		prices = append(prices, d.Price)
	}

	return prices, dates
}
