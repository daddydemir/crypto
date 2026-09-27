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
	db      *gorm.DB
	catalog *CoinCatalog
}

func NewPriceRepository(db *gorm.DB, catalog *CoinCatalog) PriceRepository {
	return &PriceRepositoryImpl{db: db, catalog: catalog}
}

func (p *PriceRepositoryImpl) GetTopCoins() ([]coincap.Coin, error) {
	catalogCoins, err := p.catalog.List()
	if err != nil {
		return nil, err
	}
	symbols := make([]string, 0, len(catalogCoins))
	for _, coin := range catalogCoins {
		symbols = append(symbols, coin.Symbol)
	}
	type price struct {
		Symbol string
		Value  float32 `gorm:"column:price_usd"`
	}
	var prices []price
	query := `select distinct on (symbol) symbol, close_price as price_usd
		from yahoo_candles where symbol in (?) and close_price is not null
		order by symbol, candle_date desc`
	if err = p.db.Raw(query, symbols).Scan(&prices).Error; err != nil {
		return nil, err
	}
	priceBySymbol := make(map[string]float32, len(prices))
	for _, current := range prices {
		priceBySymbol[current.Symbol] = current.Value
	}
	coins := make([]coincap.Coin, 0, len(catalogCoins))
	for _, coin := range catalogCoins {
		coins = append(coins, coincap.Coin{Id: coin.ID, Symbol: coin.Symbol, Name: coin.Name, PriceUsd: priceBySymbol[coin.Symbol]})
	}
	return coins, nil
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
