package cache

import "github.com/daddydemir/crypto/internal/domain/model"

type CoinCacheService interface {
	GetCoins() []model.Coin
}
