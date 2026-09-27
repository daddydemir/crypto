package infrastructure

import (
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"gorm.io/gorm"
)

type PriceRepository interface {
	GetTopCoins() ([]coincap.Coin, error)
	GetHistoricalPrices(coinID string, days int) ([]coincap.History, error)
}

type PriceRepositoryImpl struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) PriceRepository {
	return &PriceRepositoryImpl{db: db}
}

func (p *PriceRepositoryImpl) GetTopCoins() ([]coincap.Coin, error) {
	var coins []coincap.Coin
	query := `select lower(c.symbol) as id, c.symbol, coalesce(c.name, c.symbol) as name,
		coalesce(latest.close_price, 0) as price_usd
	from coins c
	left join lateral (
		select close_price from yahoo_candles
		where symbol = c.symbol and close_price is not null
		order by candle_date desc limit 1
	) latest on true
	where c.is_active
	order by c.market_cap_rank nulls last, c.symbol`
	err := p.db.Raw(query).Scan(&coins).Error
	return coins, err
}

func (p *PriceRepositoryImpl) GetHistoricalPrices(coinID string, days int) ([]coincap.History, error) {
	if days < 0 {
		days = -days
	}
	if days == 0 {
		days = 1
	}
	var list []coincap.History
	query := `select close_price as price_usd, candle_date as date
		from (
			select close_price, candle_date from yahoo_candles
			where lower(symbol) = lower(?) and close_price is not null
			order by candle_date desc limit ?
		) prices order by candle_date`
	err := p.db.Raw(query, coinID, days).Scan(&list).Error
	return list, err
}
