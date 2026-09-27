package infra

import (
	"github.com/daddydemir/crypto/pkg/channels/donchian/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRawDataWithSymbol(symbol string) ([]domain.DonchianData, error) {

	query := `
select candle_date as date, low_price as min, high_price as max, close_price as close
from yahoo_candles
where symbol = upper(?)
	and low_price is not null and high_price is not null and close_price is not null
order by candle_date asc
`
	var result []domain.DonchianData
	err := r.db.Raw(query, symbol).Scan(&result).Error
	return result, err
}
