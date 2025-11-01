package movingaverage

import (
	"context"
	"github.com/daddydemir/crypto/pkg/domain/movingaverage"
	"time"
)

type Service struct {
	PriceRepo movingaverage.PriceHistoryRepository
}

func NewService(priceRepo movingaverage.PriceHistoryRepository) *Service {
	return &Service{
		PriceRepo: priceRepo,
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
