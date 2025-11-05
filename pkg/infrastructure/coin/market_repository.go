package coin

import (
	"github.com/daddydemir/crypto/pkg/domain/coin"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"gorm.io/gorm"
)

type MarketRepository struct {
	api      *coincap.CachedClient
	database *gorm.DB
}

func NewCoinGeckoMarketRepository(api *coincap.CachedClient, db *gorm.DB) *MarketRepository {
	return &MarketRepository{api: api, database: db}
}

func (r *MarketRepository) GetCurrentPrices() ([]coin.Coin, error) {
	err, coins := r.api.ListCoins()
	if err != nil {
		return nil, err
	}
	response := make([]coin.Coin, 0, 100)
	for _, c := range coins {
		response = append(response, coin.Coin{
			ID:       c.Id,
			Name:     c.Name,
			Symbol:   c.Symbol,
			PriceUSD: float64(c.PriceUsd),
		})
	}
	return response, nil
}

func (r *MarketRepository) GetPriceChanges() ([]coin.PriceResult, error) {
	var results []coin.PriceResult
	query := `WITH price_data AS (
        SELECT 
            exchange_id,
            first_price,
            date::date as date
        FROM daily_models 
        WHERE date::date IN (CURRENT_DATE, CURRENT_DATE - INTERVAL '1 day', CURRENT_DATE - INTERVAL '7 days')
    ),
    current_prices AS (
        SELECT 
            exchange_id,
            first_price as current_price
        FROM price_data
        WHERE date::date = CURRENT_DATE
    ),
    day_ago_prices AS (
        SELECT 
            exchange_id,
            first_price as day_ago_price
        FROM price_data
        WHERE date::date = CURRENT_DATE - INTERVAL '1 day'
    ),
    week_ago_prices AS (
        SELECT 
            exchange_id,
            first_price as week_ago_price
        FROM price_data
        WHERE date::date = CURRENT_DATE - INTERVAL '7 days'
    )
    SELECT 
        c.exchange_id,
        c.current_price,
        d.day_ago_price,
        w.week_ago_price,
        ROUND(((c.current_price - d.day_ago_price) / d.day_ago_price) * 100, 2) as change_24h,
        ROUND(((c.current_price - w.week_ago_price) / w.week_ago_price) * 100, 2) as change_7d
    FROM current_prices c
    LEFT JOIN day_ago_prices d ON c.exchange_id = d.exchange_id
    LEFT JOIN week_ago_prices w ON c.exchange_id = w.exchange_id`
	err := r.database.Raw(query).Scan(&results).Error
	return results, err

}
