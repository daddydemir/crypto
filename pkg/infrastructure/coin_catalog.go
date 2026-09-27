package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/daddydemir/crypto/pkg/cache"
	"gorm.io/gorm"
)

const coinCatalogCacheKey = "crypto:coins:v1"

type Coin struct {
	ID            string `json:"id" gorm:"column:id"`
	Symbol        string `json:"symbol" gorm:"column:symbol"`
	Name          string `json:"name" gorm:"column:name"`
	MarketCapRank int    `json:"marketCapRank" gorm:"column:market_cap_rank"`
}

type CoinCatalog struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewCoinCatalog(db *gorm.DB, cacheService cache.Cache) *CoinCatalog {
	return &CoinCatalog{db: db, cache: cacheService}
}

// List reads the coin universe from Redis. On a cache miss it loads the active
// coins from PostgreSQL once and stores the complete list without expiration.
func (c *CoinCatalog) List() ([]Coin, error) {
	data, err := c.cache.Get(coinCatalogCacheKey)
	if err == nil {
		var coins []Coin
		if raw, ok := data.(string); ok && raw != "" {
			if unmarshalErr := json.Unmarshal([]byte(raw), &coins); unmarshalErr == nil && len(coins) > 0 {
				return coins, nil
			}
		}
	}

	var coins []Coin
	query := `select lower(symbol) as id, symbol, coalesce(name, symbol) as name,
		coalesce(market_cap_rank, 0) as market_cap_rank
		from coins
		where is_active
		order by market_cap_rank nulls last, symbol`
	if dbErr := c.db.Raw(query).Scan(&coins).Error; dbErr != nil {
		return nil, dbErr
	}
	if len(coins) == 0 {
		return nil, fmt.Errorf("active coin list is empty")
	}

	encoded, marshalErr := json.Marshal(coins)
	if marshalErr != nil {
		return nil, marshalErr
	}
	if cacheErr := c.cache.Set(coinCatalogCacheKey, encoded); cacheErr != nil {
		return nil, cacheErr
	}
	return coins, nil
}
