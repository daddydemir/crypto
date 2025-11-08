package exponentialma

import (
	"github.com/daddydemir/crypto/pkg/domain/exponentialma"
	"time"
)

type Service struct {
	repo exponentialma.PriceHistoryRepository
}

func NewService(repo exponentialma.PriceHistoryRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetMovingAverageSeries(coinId string, days int) ([]exponentialma.Point, error) {
	priceData, err := s.repo.GetLastNDaysPricesWithDates(coinId, days)
	if err != nil {
		return nil, err
	}
	dates, prices := make([]time.Time, 0, len(priceData)), make([]float64, 0, len(priceData))
	for _, d := range priceData {
		dates = append(dates, d.Date)
		prices = append(prices, d.Price)
	}
	return exponentialma.CalculateSeries(dates, prices), nil
}
