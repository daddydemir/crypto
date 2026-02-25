package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/bollinger/domain"
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		database: db,
	}
}

func (r *Repository) GetPrices(coinID string) ([]domain.PriceData, error) {
	sql := `select close_price as price, close_time::date as date from candles where symbol = upper(?) order by close_time`
	var results []domain.PriceData
	r.database.Raw(sql, coinID).Scan(&results)
	return results, nil
}
