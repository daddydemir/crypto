package service

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
)

type CacheService struct {
	service cache.Cache
}

func NewCacheService(c cache.Cache) *CacheService {
	return &CacheService{service: c}
}

func (c *CacheService) GetCoins() []coincap.Coin {
	data, err := c.service.Get("coinList")
	if err != nil {
		slog.Error("GetCoins", "error", err)
	}

	bytes, ok := data.(string)
	if !ok {
		slog.Error("data.(string)", "data", data)
	}

	var coins = make([]coincap.Coin, 100)
	err = json.Unmarshal([]byte(bytes), &coins)
	if err != nil {
		slog.Error("json.Unmarshall", "error", err)
	}
	return coins
}
