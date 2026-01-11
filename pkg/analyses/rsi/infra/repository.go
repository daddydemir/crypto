package infra

import (
	"errors"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/domain/coin"
	"github.com/daddydemir/crypto/pkg/domain/indicator"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/daddydemir/crypto/pkg/service"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	service *service.CacheService
	db      *gorm.DB
	cache   cache.Cache
	api     *coincap.CachedClient
}

func NewRepository(service *service.CacheService, db *gorm.DB, cache cache.Cache, api *coincap.CachedClient) *Repository {
	return &Repository{
		service: service,
		db:      db,
		cache:   cache,
		api:     api,
	}
}

func (r *Repository) GetTopCoinIDs() ([]coincap.Coin, error) {
	coins := r.service.GetCoins()
	if len(coins) == 0 {
		return nil, errors.New("coin list is empty")
	}
	return coins, nil
}

func (r *Repository) GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error) {
	before := time.Now().Add(-time.Hour * 24 * time.Duration(days))

	sql := `select lower(c.symbol) as exchange_id , c.open_time as "date", c.close_price as first_price
		from candles c 
		where lower(c.symbol) in (?)
			and c.open_time > ? 
		order by c.symbol, c.open_time`
	var results []Result
	tx := r.db.Raw(sql, ids, before.Format("2006-01-02")).Scan(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	mapp := make(map[string][]float64)
	for _, r := range results {
		mapp[r.ExchangeId] = append(mapp[r.ExchangeId], r.Price)
	}
	return mapp, nil
}

func (r *Repository) GetHistoricalPrices(coinID string, days int) ([]indicator.PriceData, error) {
	list := make([]coincap.History, 0)
	err := r.cache.GetList(coinID, &list, int64(days), -1)
	if err != nil {
		return nil, err
	}
	prices := make([]indicator.PriceData, 0, len(list))
	for _, h := range list {
		prices = append(prices, indicator.PriceData{
			Price: float64(h.PriceUsd),
			Date:  h.Date,
		})
	}

	return prices, nil
}

func (r *Repository) GetHistoricalPricesDB(coinID string) ([]indicator.PriceData, error) {
	sql := `select close_price as price, close_time::date as date from candles where symbol = upper(?) order by close_time`
	var results []indicator.PriceData
	r.db.Raw(sql, coinID).Scan(&results)
	return results, nil
}

func (r *Repository) GetPriceAt(coinId string, date int) (float64, error) {
	var price float64
	list := make([]coincap.History, 0)
	err := r.cache.GetList(coinId, &list, 0, -1)
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

func (r *Repository) GetCurrentPrices() ([]coin.Coin, error) {
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

func (r *Repository) GetPriceChanges() ([]coin.PriceResult, error) {
	var results []coin.PriceResult
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

type Result struct {
	ExchangeId string
	Date       string
	Price      float64 `gorm:"column:first_price"`
}
