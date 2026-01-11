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
select open_time::date as date, low_price as min, high_price as max, close_price as close
from candles
where symbol = upper(?)
order by open_time asc
`
	var result []domain.DonchianData
	err := r.db.Raw(query, symbol).Scan(&result).Error
	return result, err
}
