package coin

import (
	"errors"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"time"
)

type CacheHistoryRepository struct {
	cacheService cache.Cache
}

func NewCacheHistoryRepository(cacheService cache.Cache) *CacheHistoryRepository {
	return &CacheHistoryRepository{
		cacheService: cacheService,
	}
}

func (c *CacheHistoryRepository) GetPriceAt(coinId string, date int) (float64, error) {
	var price float64
	list := make([]coincap.History, 0)
	err := c.cacheService.GetList(coinId, &list, 0, -1)
	if err != nil {
		return price, err
	}
	if len(list) < date {
		return price, errors.New("list size is not valid")
	}
	if list[len(list)-1].Date.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		price = float64(list[len(list)-date-1].PriceUsd)
	} else {
		price = float64(list[len(list)-date].PriceUsd)
	}

	return price, nil
}
