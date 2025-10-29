package coin

import (
	"errors"
	"github.com/daddydemir/crypto/pkg/cache"
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

func (p *PriceRepository) GetTopCoinIDs() ([]string, error) {
	coins := p.service.GetCoins()
	if len(coins) == 0 {
		return nil, errors.New("coin list is empty")
	}
	response := make([]string, 0, len(coins))
	for _, c := range coins {
		response = append(response, c.Id)
	}
	return response, nil
}

func (p *PriceRepository) GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error) {
	before := time.Now().Add(-time.Hour * 24 * time.Duration(days))

	sql := `select lower(em.name) as exchange_id , dm.date, dm.first_price 
		from daily_models dm, exchange_models em 
		where lower(em.name) in (?)
			and em.exchange_id = dm.exchange_id
			and dm.date > ? 
		order by em.name, dm.date`
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
