package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/coin/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCurrentPrices() ([]domain.Coin, error) {
	var coins []domain.Coin
	query := `select lower(c.symbol) as id, coalesce(c.name, c.symbol) as name,
		c.symbol, coalesce(latest.close_price, 0) as price_usd
	from coins c
	left join lateral (
		select yc.close_price
		from yahoo_candles yc
		where yc.symbol = c.symbol and yc.close_price is not null
		order by yc.candle_date desc
		limit 1
	) latest on true
	where c.is_active
	order by c.market_cap_rank nulls last, c.symbol`
	err := r.db.Raw(query).Scan(&coins).Error
	return coins, err
}

func (r *Repository) GetPriceChanges() ([]domain.PriceResult, error) {
	var results []domain.PriceResult
	query := `
with latest as (
	select distinct on (symbol) symbol, candle_date, close_price as current_price
	from yahoo_candles
	where close_price is not null
	order by symbol, candle_date desc
)
select l.symbol as exchange_id,
	l.current_price,
	coalesce(d.close_price, 0) as day_ago_price,
	coalesce(w.close_price, 0) as week_ago_price,
	coalesce(m.close_price, 0) as month_ago_price,
	coalesce(a.week_avg_price, 0) as avg_7_days_price,
	coalesce(a.month_avg_price, 0) as avg_30_days_price,
	coalesce(round(((l.current_price - d.close_price) / nullif(d.close_price, 0)) * 100, 2), 0) as change_24h,
	coalesce(round(((l.current_price - w.close_price) / nullif(w.close_price, 0)) * 100, 2), 0) as change_7d,
	coalesce(round(((l.current_price - m.close_price) / nullif(m.close_price, 0)) * 100, 2), 0) as change_30d,
	coalesce(round(((l.current_price - a.week_avg_price) / nullif(a.week_avg_price, 0)) * 100, 2), 0) as change_arithmetic_7d,
	coalesce(round(((l.current_price - a.month_avg_price) / nullif(a.month_avg_price, 0)) * 100, 2), 0) as change_arithmetic_30d
from latest l
join coins c on c.symbol = l.symbol and c.is_active
left join lateral (select close_price from yahoo_candles where symbol = l.symbol and close_price is not null and candle_date <= l.candle_date - 1 order by candle_date desc limit 1) d on true
left join lateral (select close_price from yahoo_candles where symbol = l.symbol and close_price is not null and candle_date <= l.candle_date - 7 order by candle_date desc limit 1) w on true
left join lateral (select close_price from yahoo_candles where symbol = l.symbol and close_price is not null and candle_date <= l.candle_date - 30 order by candle_date desc limit 1) m on true
left join lateral (
	select avg(close_price) filter (where candle_date > l.candle_date - 7) as week_avg_price,
		avg(close_price) filter (where candle_date > l.candle_date - 30) as month_avg_price
	from yahoo_candles
	where symbol = l.symbol and close_price is not null and candle_date < l.candle_date
) a on true
`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err

}
