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

	query := `select c.candle_date as date
	, c.open_price as open
	, c.high_price as high
	, c.low_price as low
	, c.close_price as close
	, c.volume
from yahoo_candles c
where c.symbol = upper(?)
	and c.open_price is not null and c.high_price is not null and c.low_price is not null and c.close_price is not null
order by c.candle_date`

	var result []domain.PriceData
	err := r.db.Raw(query, symbol).Scan(&result).Error
	return result, err
}
