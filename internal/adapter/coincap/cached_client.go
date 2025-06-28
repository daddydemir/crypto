package coincapadapter

import (
	"fmt"
	"time"

	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
)

type CachedCoinCapClient struct {
	inner        port.CoinCapAPI
	cacheService port.Cache
}

func NewCachedCoinCapClient(inner port.CoinCapAPI, cache port.Cache) port.CoinCapAPI {
	return &CachedCoinCapClient{
		inner:        inner,
		cacheService: cache,
	}
}

func (c *CachedCoinCapClient) ListCoins() (error, []model.Coin) {
	key := "listCoins"
	var data []model.Coin

	err := c.cacheService.GetList(key, &data, 0, -1)
	if err == nil && len(data) > 0 {
		return nil, data
	}

	err, data = c.inner.ListCoins()
	if err != nil {
		return err, data
	}

	_ = c.cacheService.SetList(key, data, 4*time.Hour)
	return nil, data
}

func (c *CachedCoinCapClient) HistoryWithId(id string) (error, []model.History) {
	key := "history:" + id
	var data []model.History

	err := c.cacheService.GetList(key, &data, 0, -1)
	if err == nil && len(data) > 0 {
		return nil, data
	}

	err, data = c.inner.HistoryWithId(id)
	if err != nil {
		return err, data
	}

	_ = c.cacheService.SetList(key, data, 4*time.Hour)
	return nil, data
}

func (c *CachedCoinCapClient) HistoryWithTime(id string, start, end int64) (error, []model.History) {
	key := fmt.Sprintf("history:%s-%d-%d", id, start, end)
	var data []model.History

	err := c.cacheService.GetList(key, &data, 0, -1)
	if err == nil && len(data) > 0 {
		return nil, data
	}

	err, data = c.inner.HistoryWithTime(id, start, end)
	if err != nil {
		return err, data
	}

	_ = c.cacheService.SetList(key, data, 4*time.Hour)
	return nil, data
}
