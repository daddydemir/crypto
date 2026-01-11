package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/adi/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetRawDataWithSymbol(symbol string) ([]domain.PriceData, error) {

	query := `select c.open_time as date
	, c.open_price as open
	, c.high_price as high
	, c.low_price as low
	, c.close_price as close
	, c.volume
from candles c 
where c.symbol = upper(?)
order by c.open_time`

	var result []domain.PriceData
	err := r.db.Raw(query, symbol).Scan(&result).Error
	return result, err
}
