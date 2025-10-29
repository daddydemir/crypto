package coincap

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/pkg/cache"
	"log/slog"
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

	err = c.saveToBothCaches(key, coins)
	if err != nil {
		slog.Error("Failed to save to cache", "error", err)
	}

	return nil, coins
}

func (c *CachedClient) saveToBothCaches(key string, coins []Coin) error {
	err := c.cacheService.SetList(key, coins, time.Hour*4)
	if err != nil {
		return fmt.Errorf("failed to save list to cache: %w", err)
	}

	jsonData, err := json.Marshal(coins)
	if err != nil {
		return fmt.Errorf("failed to marshal coins to JSON: %w", err)
	}

	err = c.cacheService.Set("coinList", jsonData)
	if err != nil {
		return fmt.Errorf("failed to save JSON to cache: %w", err)
	}

	return nil
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
