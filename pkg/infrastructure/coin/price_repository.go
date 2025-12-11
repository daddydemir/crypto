package coin

import (
	"errors"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/domain/indicator"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/daddydemir/crypto/pkg/service"
	"gorm.io/gorm"
	"time"
)

type PriceRepository struct {
	cacheService cache.Cache
	service      *service.CacheService
	database     *gorm.DB
}

func NewPriceRepository(cacheService cache.Cache, service *service.CacheService, database *gorm.DB) *PriceRepository {
	return &PriceRepository{
		cacheService: cacheService,
		service:      service,
		database:     database,
	}
}

func (p *PriceRepository) GetTopCoinIDs() ([]coincap.Coin, error) {
	coins := p.service.GetCoins()
	if len(coins) == 0 {
		return nil, errors.New("coin list is empty")
	}
	return coins, nil
}

func (p *PriceRepository) GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error) {
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

func (p *PriceRepository) GetHistoricalPrices(coinID string, days int) ([]indicator.PriceData, error) {
	list := make([]coincap.History, 0)
	err := p.cacheService.GetList(coinID, &list, int64(days), -1)
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

func (p *PriceRepository) GetHistoricalPricesDB(coinID string) ([]indicator.PriceData, error) {
	sql := `select close_price as price, close_time::date as date from candles where symbol = upper(?) order by close_time`
	var results []indicator.PriceData
	p.database.Raw(sql, coinID).Scan(&results)
	return results, nil
}
