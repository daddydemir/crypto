package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/macd/domain"
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
}

type Result struct {
	ExchangeId string
	Date       string
	Price      float64 `gorm:"column:first_price"`
}

func NewRepository(database *gorm.DB) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) GetCoinPricesFromDB(coinID string) ([]domain.PriceData, error) {

	sql := `select close_price as price, close_time as date 
			from candles 
			where lower(symbol) = lower(?) 
			order by close_time`

	var results []domain.PriceData
	tx := r.database.Raw(sql, coinID).Scan(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return results, nil
}
