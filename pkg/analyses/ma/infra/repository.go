package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/ma/domain"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
)

type Repository struct {
	cacheService cache.Cache
}

func NewRepository(cacheService cache.Cache) *Repository {
	return &Repository{
		cacheService: cacheService,
	}
}

func (r *Repository) GetLastNDaysPricesWithDates(coinID string, days int) ([]domain.PriceData, error) {
	list := make([]coincap.History, 0)
	err := r.cacheService.GetList(coinID, &list, 0, int64(days))
	if err != nil {
		return nil, err
	}
	prices := make([]domain.PriceData, 0, len(list))
	for _, h := range list {
		prices = append(prices, domain.PriceData{
			Price: float64(h.PriceUsd),
			Date:  h.Date,
		})
	}
	return prices, nil
}
