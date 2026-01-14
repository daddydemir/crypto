package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/coin/domain"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"gorm.io/gorm"
)

type Repository struct {
	api *coincap.CachedClient
	db  *gorm.DB
}

func NewRepository(api *coincap.CachedClient, db *gorm.DB) *Repository {
	return &Repository{api: api, db: db}
}

func (r *Repository) GetCurrentPrices() ([]domain.Coin, error) {
	err, coins := r.api.ListCoins()
	if err != nil {
		return nil, err
	}
	response := make([]domain.Coin, 0, 100)
	for _, c := range coins {
		response = append(response, domain.Coin{
			ID:       c.Id,
			Name:     c.Name,
			Symbol:   c.Symbol,
			PriceUSD: float64(c.PriceUsd),
		})
	}
	return response, nil
}

func (r *Repository) GetPriceChanges() ([]domain.PriceResult, error) {
	var results []domain.PriceResult
	query := `WITH price_data AS (
        SELECT 
            exchange_id,
            first_price,
            date::date as date
        FROM daily_models 
        WHERE date::date IN (CURRENT_DATE, CURRENT_DATE - INTERVAL '1 day', CURRENT_DATE - INTERVAL '7 days', CURRENT_DATE - INTERVAL '30 days')
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
    ),
	month_ago_prices AS (
	SELECT 
		exchange_id,
		first_price month_ago_price
	FROM price_data
	WHERE date::date = CURRENT_DATE - INTERVAL '30 days'
	),
	avg_7_days_price AS (
		SELECT
			exchange_id,
			AVG(first_price) as avg_price
		FROM daily_models
		WHERE date::date BETWEEN CURRENT_DATE - INTERVAL '7 days' AND CURRENT_DATE - INTERVAL '1 day'
		GROUP BY exchange_id
	),
	avg_30_days_price AS (
		SELECT
			exchange_id,
			avg(first_price) as avg_price
		FROM daily_models
		WHERE date::date BETWEEN CURRENT_DATE - INTERVAL '30 days' AND CURRENT_DATE - INTERVAL '1 day'
		GROUP BY exchange_id
	)
    SELECT 
        c.exchange_id,
        c.current_price,
        d.day_ago_price,
        w.week_ago_price,
		m.month_ago_price,
		avg_7.avg_price as avg_7_days_price,
		avg_30.avg_price as avg_30_days_price,
        ROUND(((c.current_price - d.day_ago_price) / d.day_ago_price) * 100, 2) as change_24h,
        ROUND(((c.current_price - w.week_ago_price) / w.week_ago_price) * 100, 2) as change_7d,
		ROUND(((c.current_price - m.month_ago_price) / m.month_ago_price) * 100, 2) as change_30d,
		ROUND(((c.current_price - avg_7.avg_price) /  avg_7.avg_price) * 100 , 2) as change_arithmetic_7d,
		ROUND(((c.current_price - avg_30.avg_price) / avg_30.avg_price) * 100 , 2) as change_arithmetic_30d
    FROM current_prices c
    LEFT JOIN day_ago_prices d ON c.exchange_id = d.exchange_id
    LEFT JOIN week_ago_prices w ON c.exchange_id = w.exchange_id
	LEFT JOIN month_ago_prices m ON c.exchange_id = m.exchange_id
	LEFT JOIN avg_7_days_price avg_7 ON c.exchange_id = avg_7.exchange_id
	LEFT JOIN avg_30_days_price avg_30 ON c.exchange_id = avg_30.exchange_id`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err

}
