package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/rsi/domain"
	"github.com/daddydemir/crypto/pkg/infrastructure"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"gorm.io/gorm"
	"time"
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

func (p *Repository) GetTopCoinIDs() ([]coincap.Coin, error) {
	catalogCoins, err := p.catalog.List()
	if err != nil {
		return nil, err
	}
	coins := make([]coincap.Coin, 0, len(catalogCoins))
	for _, coin := range catalogCoins {
		coins = append(coins, coincap.Coin{Id: coin.ID, Symbol: coin.Symbol, Name: coin.Name})
	}
	return coins, nil
}

func (p *Repository) GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error) {
	before := time.Now().Add(-time.Hour * 24 * time.Duration(days))

	sql := `select lower(c.symbol) as exchange_id, c.candle_date as "date", c.close_price as first_price
		from yahoo_candles c
		where lower(c.symbol) in (?)
			and c.candle_date > ? and c.close_price is not null
		order by c.symbol, c.candle_date`
	var results []Result
	tx := p.database.Raw(sql, ids, before.Format("2006-01-02")).Scan(&results)
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
