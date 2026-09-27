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
	query := `select symbol, candle_date as time, close_price as close, high_price as high, low_price as low
	from yahoo_candles where symbol = upper(?)
	and close_price is not null and high_price is not null and low_price is not null order by candle_date`
	var result []domain.Candle
	err := r.db.Raw(query, symbol).Scan(&result).Error
	return result, err
}

func (r *CandleRepository) GetBySymbolAndYear(symbol, year string) ([]domain.Candle, error) {
	query := `select symbol, candle_date as time, close_price as close
		, high_price as high, low_price as low
		from yahoo_candles
		where symbol = upper(?) 
		and extract(year from candle_date)::text = ?
		and close_price is not null and high_price is not null and low_price is not null
		order by candle_date`
	var result []domain.Candle
	err := r.db.Raw(query, symbol, year).Scan(&result).Error
	return result, err
}

func (r *CandleRepository) GetBySymbolAndYearMonth(symbol, year, month string) ([]domain.Candle, error) {
	query := `select symbol, candle_date as time, close_price as close
		, high_price as high, low_price as low
		from yahoo_candles
		where symbol = upper(?) 
		and extract(year from candle_date)::text = ?
		and lpad(extract(month from candle_date)::text, 2, '0') = ?
		and close_price is not null and high_price is not null and low_price is not null
		order by candle_date`
	var result []domain.Candle
	err := r.db.Raw(query, symbol, year, month).Scan(&result).Error
	return result, err
}
