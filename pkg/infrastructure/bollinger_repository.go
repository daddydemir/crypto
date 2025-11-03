package infrastructure

import (
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/domain/bollinger"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
)

type BollingerRepository struct {
	cacheService cache.Cache
}

func NewBollingerRepository(cacheService cache.Cache) *BollingerRepository {
	return &BollingerRepository{
		cacheService: cacheService,
	}
}

func (b *BollingerRepository) GetLastNDaysPricesWithDates(coinID string, days int) ([]bollinger.PriceData, error) {
	list := make([]coincap.History, 0)
	err := b.cacheService.GetList(coinID, &list, 0, int64(days))
	if err != nil {
		return nil, err
	}
	prices := make([]bollinger.PriceData, 0, len(list))
	for _, h := range list {
		prices = append(prices, bollinger.PriceData{
			Price: float64(h.PriceUsd),
			Date:  h.Date,
		})
	}
	return prices, nil
}
