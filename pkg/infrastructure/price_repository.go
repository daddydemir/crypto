package infrastructure

import (
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/daddydemir/crypto/pkg/service"
)

type PriceRepository interface {
	GetTopCoins() ([]coincap.Coin, error)
	GetHistoricalPrices(coinID string, days int) ([]coincap.History, error)
}

type PriceRepositoryImpl struct {
	service *service.CacheService
	cache   cache.Cache
}

func NewPriceRepository(service *service.CacheService, cache cache.Cache) PriceRepository {
	return &PriceRepositoryImpl{
		service: service,
		cache:   cache,
	}
}

func (p *PriceRepositoryImpl) GetTopCoins() ([]coincap.Coin, error) {
	return p.service.GetCoins(), nil
}

func (p *PriceRepositoryImpl) GetHistoricalPrices(coinID string, days int) ([]coincap.History, error) {
	list := make([]coincap.History, 0)
	err := p.cache.GetList(coinID, &list, int64(days), -1)
	if err != nil {
		return nil, err
	}
	return list, nil
}
