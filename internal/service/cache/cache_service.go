package cache

import (
	"encoding/json"
	"github.com/daddydemir/crypto/internal/port"
	"log/slog"

	"github.com/daddydemir/crypto/internal/domain/model"
)

type CacheService struct {
	cache port.Cache
}

func NewCacheService(c port.Cache) *CacheService {
	return &CacheService{cache: c}
}

func (c *CacheService) GetCoins() []model.Coin {
	data, err := c.cache.Get("coinList")
	if err != nil {
		slog.Error("CacheService:Get", "error", err)
		return nil
	}

	bytes, ok := data.(string)
	if !ok {
		slog.Error("CacheService:assert string", "data", data)
		return nil
	}

	var coins []model.Coin
	err = json.Unmarshal([]byte(bytes), &coins)
	if err != nil {
		slog.Error("CacheService:unmarshal", "error", err)
		return nil
	}

	return coins
}
