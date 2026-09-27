package infra

import (
	"github.com/daddydemir/crypto/pkg/analyses/atr/domain"
	"gorm.io/gorm"
	"strings"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetPointsBySymbol(symbol string) ([]domain.AtrPoint, error) {
	query := `
select symbol, high_price as current_high, low_price as current_low,
	previous_close as yesterday_close, candle_date as time
from (
	select symbol, candle_date, high_price, low_price,
		lag(close_price) over (partition by symbol order by candle_date) as previous_close
	from yahoo_candles
	where symbol = ? and high_price is not null and low_price is not null and close_price is not null
) prices
where previous_close is not null
order by candle_date`

	var result []domain.AtrPoint
	err := r.db.Raw(query, strings.ToUpper(symbol)).Scan(&result).Error
	return result, err
}
