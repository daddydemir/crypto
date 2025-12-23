package infrastructure

import (
	"github.com/daddydemir/crypto/pkg/atr/domain"
	"gorm.io/gorm"
	"strings"
)

type AtrRepository struct {
	db *gorm.DB
}

func NewAtrRepository(db *gorm.DB) *AtrRepository {
	return &AtrRepository{
		db: db,
	}
}

func (r *AtrRepository) GetPointsBySymbol(symbol string) ([]domain.AtrPoint, error) {
	query := `
select c.symbol
	, c.high_price as current_high
	, c.low_price as current_low
	, c.open_price as yesterday_close
	, c.open_time::date as time
from candles c 
where c.symbol = ?
order by c.open_time`

	var result []domain.AtrPoint
	err := r.db.Raw(query, strings.ToUpper(symbol)).Scan(&result).Error
	return result, err
}
