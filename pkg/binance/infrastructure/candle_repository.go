package infrastructure

import (
	"github.com/daddydemir/crypto/pkg/binance/domain"
	"gorm.io/gorm"
)

type CandleRepository struct {
	db *gorm.DB
}

func NewCandleRepository(db *gorm.DB) *CandleRepository {
	return &CandleRepository{db: db}
}

func (r *CandleRepository) GetBySymbol(symbol string) ([]domain.Candle, error) {
	query := `select symbol, close_time as time, close_price as close from candles where symbol = upper(?) order by close_time`
	var result []domain.Candle
	err := r.db.Raw(query, symbol).Scan(&result).Error
	return result, err
}
