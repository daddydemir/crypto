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
	query := `
with price_data as (
	select c.symbol, c.close_price, c.close_time::date
	from candles c 
	where c.close_time::date in (current_date - interval '1 day', current_date - interval '2 days', current_date - interval '8 days', current_date - interval '31 days')
),
current_price as (
	select symbol, close_price as current_price
	from price_data
	where close_time::date = current_date - interval '1 day'
), day_ago_prices as (
	select symbol, close_price as day_ago_price
	from price_data
	where close_time::date = current_date - interval '2 days'
), week_ago_prices as (
	select symbol, close_price as week_ago_price
	from price_data 
	where close_time::date = current_date - interval '8 days'
), month_ago_prices as (
	select symbol, close_price as month_ago_price
	from price_data
	where close_time::date = current_date - interval '31 days'
), avg_7_days_price as (
	select symbol, avg(close_price) as week_avg_price
	from candles c
	where c.open_time::date between current_date - interval '8 days' and current_date - interval '2 days'
	group by symbol
), avg_30_days_price as (
	select symbol, avg(close_price) as month_avg_price
	from candles c
	where c.open_time::date between current_date - interval '31 days' and current_date - interval '2 days'
	group by symbol
)
select cp.symbol as exchange_id
	, cp.current_price
	, dap.day_ago_price
	, wap.week_ago_price
	, map.month_ago_price
	, a7.week_avg_price as avg_7_days_price
	, a30.month_avg_price as avg_30_days_price
	, round(((cp.current_price - dap.day_ago_price) / dap.day_ago_price) * 100, 2) as change_24h
	, round(((cp.current_price - wap.week_ago_price) / wap.week_ago_price) * 100, 2) as change_7d
	, round(((cp.current_price - map.month_ago_price) / map.month_ago_price) * 100, 2) as change_30d
	, round(((cp.current_price - a7.week_avg_price) / a7.week_avg_price) * 100, 2) as change_arithmetic_7d
	, round(((cp.current_price - a30.month_avg_price) / a30.month_avg_price) * 100, 2) as change_arithmetic_30d
from current_price cp
	left join day_ago_prices dap on dap.symbol = cp.symbol
	left join week_ago_prices wap on wap.symbol = cp.symbol
	left join month_ago_prices map on map.symbol = cp.symbol
	left join avg_7_days_price a7 on a7.symbol = cp.symbol
	left join avg_30_days_price a30 on a30.symbol = cp.symbol
`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err

}
