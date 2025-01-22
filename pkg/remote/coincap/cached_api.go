package coincap

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/cache"
	"time"
)

type CachedClient struct {
	originalClient Client
	cacheService   cache.Cache
}

func NewCachedClient(client Client, service cache.Cache) *CachedClient {
	return &CachedClient{
		originalClient: client,
		cacheService:   service,
	}
}

func (c *CachedClient) ListCoins() (error, []Coin) {
	key := "listCoins"
	var coins []Coin
	var err error
	err = c.cacheService.GetList(key, &coins, 0, -1)
	if coins != nil {
		return err, coins
	}
	err, coins = c.originalClient.ListCoins()
	if err != nil {
		return err, coins
	}
	err = c.cacheService.SetList(key, coins, time.Hour*4)
	return nil, coins
}

func (c *CachedClient) HistoryWithId(s string) (error, []History) {
	key := "history:" + s
	var histories []History
	var err error

	err = c.cacheService.GetList(key, &histories, 0, -1)
	if histories != nil {
		return err, histories
	}

	err, histories = c.originalClient.HistoryWithId(s)
	if err != nil {
		return err, histories
	}

	err = c.cacheService.SetList(key, histories, time.Hour*4)
	return err, histories
}

func (c *CachedClient) HistoryWithTime(s string, start, end int64) (error, []History) {
	key := fmt.Sprintf("history:%s-%d-%d", s, start, end)
	var histories []History
	var err error

	err = c.cacheService.GetList(key, &histories, 0, -1)
	if histories != nil {
		return err, histories
	}

	err, histories = c.originalClient.HistoryWithTime(s, start, end)
	if err != nil {
		return err, histories
	}

	err = c.cacheService.SetList(key, histories, time.Hour*4)
	return err, histories
}
