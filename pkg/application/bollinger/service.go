package bollinger

import (
	"context"
	"github.com/daddydemir/crypto/pkg/domain/bollinger"

	"time"
)

type Service struct {
	PriceRepo bollinger.PriceHistoryRepository
}

func NewService(priceRepo bollinger.PriceHistoryRepository) *Service {
	return &Service{
		PriceRepo: priceRepo,
	}
}

func (s *Service) GetBollingerSeries(_ context.Context, coinID string, days int) ([]bollinger.Point, error) {
	datas, err := s.PriceRepo.GetLastNDaysPricesWithDates(coinID, days)
	if err != nil {
		return nil, err
	}
	dates, prices := make([]time.Time, 0, len(datas)), make([]float64, 0, len(datas))
	for _, d := range datas {
		dates = append(dates, d.Date)
		prices = append(prices, d.Price)
	}
	return bollinger.CalculateBollinger(prices, dates), nil
}
