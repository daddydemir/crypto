package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/rsi/domain"
	"github.com/daddydemir/crypto/pkg/infrastructure"
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
	catalog  *infrastructure.CoinCatalog
}
type Result struct {
	ExchangeId string
	Date       string
	Price      float64 `gorm:"column:first_price"`
}

func NewRepository(database *gorm.DB, catalog *infrastructure.CoinCatalog) *Repository {
	return &Repository{database: database, catalog: catalog}
}

func (p *Repository) GetTopCoinIDs() ([]domain.Coin, error) {
	catalogCoins, err := p.catalog.List()
	if err != nil {
		return nil, err
	}
	coins := make([]domain.Coin, 0, len(catalogCoins))
	for _, coin := range catalogCoins {
		coins = append(coins, domain.Coin{Symbol: coin.Symbol, Name: coin.Name})
	}
	return coins, nil
}

func (p *Repository) GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error) {
	sql := `select exchange_id, "date", first_price
		from (
			select lower(c.symbol) as exchange_id,
				c.candle_date as "date",
				c.close_price as first_price,
				row_number() over (partition by lower(c.symbol) order by c.candle_date desc) as row_number
			from yahoo_candles c
			where lower(c.symbol) in (?) and c.close_price is not null
		) ranked_prices
		where row_number <= ?
		order by exchange_id, "date"`
	var results []Result
	tx := p.database.Raw(sql, ids, days).Scan(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	mapp := make(map[string][]float64)
	for _, r := range results {
		mapp[r.ExchangeId] = append(mapp[r.ExchangeId], r.Price)
	}
	return mapp, nil
}

func (p *Repository) GetHistoricalPricesDB(coinID string) ([]domain.PriceData, error) {
	sql := `select close_price as price, candle_date as date from yahoo_candles where symbol = upper(?) and close_price is not null order by candle_date`
	var results []domain.PriceData
	p.database.Raw(sql, coinID).Scan(&results)
	return results, nil
}
