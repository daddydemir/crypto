package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/coin/domain"
	"github.com/daddydemir/crypto/pkg/infrastructure"
	"gorm.io/gorm"
)

type Repository struct {
	db      *gorm.DB
	catalog *infrastructure.CoinCatalog
}

func NewRepository(db *gorm.DB, catalog *infrastructure.CoinCatalog) *Repository {
	return &Repository{db: db, catalog: catalog}
}

func (r *Repository) GetCurrentPrices() ([]domain.Coin, error) {
	catalogCoins, err := r.catalog.List()
	if err != nil {
		return nil, err
	}
	symbols := coinSymbols(catalogCoins)
	type price struct {
		Symbol string
		Value  float64 `gorm:"column:price_usd"`
	}
	var prices []price
	query := `select distinct on (symbol) symbol, close_price as price_usd
		from yahoo_candles where symbol in (?) and close_price is not null
		order by symbol, candle_date desc`
	if err = r.db.Raw(query, symbols).Scan(&prices).Error; err != nil {
		return nil, err
	}
	priceBySymbol := make(map[string]float64, len(prices))
	for _, current := range prices {
		priceBySymbol[current.Symbol] = current.Value
	}
	coins := make([]domain.Coin, 0, len(catalogCoins))
	for _, coin := range catalogCoins {
		coins = append(coins, domain.Coin{ID: coin.ID, Name: coin.Name, Symbol: coin.Symbol, PriceUSD: priceBySymbol[coin.Symbol]})
	}
	return coins, nil
}

func (r *Repository) GetPriceChanges() ([]domain.PriceResult, error) {
	coins, err := r.catalog.List()
	if err != nil {
		return nil, err
	}
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
left join lateral (select close_price from yahoo_candles where symbol = l.symbol and close_price is not null and candle_date <= l.candle_date - 1 order by candle_date desc limit 1) d on true
left join lateral (select close_price from yahoo_candles where symbol = l.symbol and close_price is not null and candle_date <= l.candle_date - 7 order by candle_date desc limit 1) w on true
left join lateral (select close_price from yahoo_candles where symbol = l.symbol and close_price is not null and candle_date <= l.candle_date - 30 order by candle_date desc limit 1) m on true
left join lateral (
	select avg(close_price) filter (where candle_date > l.candle_date - 7) as week_avg_price,
		avg(close_price) filter (where candle_date > l.candle_date - 30) as month_avg_price
	from yahoo_candles
	where symbol = l.symbol and close_price is not null and candle_date < l.candle_date
) a on true
where l.symbol in (?)
`
	err = r.db.Raw(query, coinSymbols(coins)).Scan(&results).Error
	return results, err

}

func coinSymbols(coins []infrastructure.Coin) []string {
	symbols := make([]string, 0, len(coins))
	for _, coin := range coins {
		symbols = append(symbols, coin.Symbol)
	}
	return symbols
}
