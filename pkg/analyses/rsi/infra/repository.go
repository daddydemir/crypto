package infra

import (
	"errors"
	"github.com/daddydemir/crypto/pkg/analyses/rsi/domain"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/daddydemir/crypto/pkg/service"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	cacheService cache.Cache
	service      *service.CacheService
	database     *gorm.DB
}
type Result struct {
	ExchangeId string
	Date       string
	Price      float64 `gorm:"column:first_price"`
}

func NewRepository(cacheService cache.Cache, service *service.CacheService, database *gorm.DB) *Repository {
	return &Repository{
		cacheService: cacheService,
		service:      service,
		database:     database,
	}
}

func (p *Repository) GetTopCoinIDs() ([]coincap.Coin, error) {
	coins := p.service.GetCoins()
	if len(coins) == 0 {
		return nil, errors.New("coin list is empty")
	}
	return coins, nil
}

func (p *Repository) GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error) {
	before := time.Now().Add(-time.Hour * 24 * time.Duration(days))

	sql := `select lower(c.symbol) as exchange_id , c.open_time as "date", c.close_price as first_price
		from candles c 
		where lower(c.symbol) in (?)
			and c.open_time > ? 
		order by c.symbol, c.open_time`
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

func (p *Repository) GetHistoricalPrices(coinID string, days int) ([]domain.PriceData, error) {
	list := make([]coincap.History, 0)
	err := p.cacheService.GetList(coinID, &list, int64(days), -1)
	if err != nil {
		return nil, err
	}
	prices := make([]domain.PriceData, 0, len(list))
	for _, h := range list {
		prices = append(prices, domain.PriceData{
			Price: float64(h.PriceUsd),
			Date:  h.Date,
		})
	}

	return prices, nil
}

func (p *Repository) GetHistoricalPricesDB(coinID string) ([]domain.PriceData, error) {
	sql := `select close_price as price, close_time::date as date from candles where symbol = upper(?) order by close_time`
	var results []domain.PriceData
	p.database.Raw(sql, coinID).Scan(&results)
	return results, nil
}
