package infrastructure

import (
	"context"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/domain/movingaverage"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
)

type PriceHistoryRepository struct {
	cacheService cache.Cache
}

func NewPriceHistoryRepository(cacheService cache.Cache) *PriceHistoryRepository {
	return &PriceHistoryRepository{
		cacheService: cacheService,
	}
}

func (p *PriceHistoryRepository) GetLastNDaysPricesWithDates(ctx context.Context, coinID string, days int) ([]movingaverage.PriceData, error) {
	list := make([]coincap.History, 0)
	err := p.cacheService.GetList(coinID, &list, 0, int64(days))
	if err != nil {
		return nil, err
	}
	prices := make([]movingaverage.PriceData, 0, len(list))
	for _, h := range list {
		prices = append(prices, movingaverage.PriceData{
			Price: float64(h.PriceUsd),
			Date:  h.Date,
		})
	}
	return prices, nil
}
