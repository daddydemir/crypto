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
	query := `select symbol, (close_time - interval '3 hours')::date as time, close_price as close from candles where symbol = upper(?) order by close_time`
	var result []domain.Candle
	err := r.db.Raw(query, symbol).Scan(&result).Error
	return result, err
}

func (r *CandleRepository) GetBySymbolAndYear(symbol, year string) ([]domain.Candle, error) {
	query := `select symbol, (close_time - interval '3 hours')::date as time, close_price as close 
		from candles 
		where symbol = upper(?) 
		and to_char(((close_time - interval '3 hours')::date), 'YYYY') = ?
		order by close_time`
	var result []domain.Candle
	err := r.db.Raw(query, symbol, year).Scan(&result).Error
	return result, err
}

func (r *CandleRepository) GetBySymbolAndYearMonth(symbol, year, month string) ([]domain.Candle, error) {
	query := `select symbol, (close_time - interval '3 hours')::date as time, close_price as close 
		from candles 
		where symbol = upper(?) 
		and to_char(((close_time - interval '3 hours')::date), 'YYYY') = ?
		and to_char(((close_time - interval '3 hours')::date), 'MM') = ?
		order by close_time`
	var result []domain.Candle
	err := r.db.Raw(query, symbol, year, month).Scan(&result).Error
	return result, err
}
