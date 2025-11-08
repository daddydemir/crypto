package exponentialma

import (
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/domain/exponentialma"
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

func (p *PriceHistoryRepository) GetLastNDaysPricesWithDates(coinID string, days int) ([]exponentialma.PriceData, error) {
	list := make([]coincap.History, 0)
	err := p.cacheService.GetList(coinID, &list, 0, int64(days))
	if err != nil {
		return nil, err
	}
	prices := make([]exponentialma.PriceData, 0, len(list))
	for _, h := range list {
		prices = append(prices, exponentialma.PriceData{
			Price: float64(h.PriceUsd),
			Date:  h.Date,
		})
	}
	return prices, nil
}
