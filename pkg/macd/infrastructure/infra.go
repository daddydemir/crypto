package infrastructure

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/daddydemir/crypto/pkg/macd/domain"
	"gorm.io/gorm"
)

// DatabasePriceRepository implements PriceRepository using the existing daily table
type DatabasePriceRepository struct {
	db *gorm.DB
}

func NewDatabasePriceRepository(db *gorm.DB) *DatabasePriceRepository {
	return &DatabasePriceRepository{db: db}
}

func (p *DatabasePriceRepository) GetPriceData(symbol string, from, to time.Time) ([]domain.PriceData, error) {
	slog.Info("Getting price data from candles table", "symbol", symbol, "from", from, "to", to)

	var prices []domain.PriceData

	query := `select open_time as date, close_price as price
from candles 
where symbol = upper(?) and close_time between ? and ?
order by close_time`

	p.db.Raw(query, symbol, from, to).Scan(&prices)

	if len(prices) == 0 {
		return nil, fmt.Errorf("no valid price data found for symbol: %s", symbol)
	}

	slog.Info("Daily price data retrieved", "symbol", symbol, "count", len(prices))
	return prices, nil
}

func (p *DatabasePriceRepository) GetLatestPriceData(symbol string, limit int) ([]domain.PriceData, error) {
	slog.Info("Getting latest daily price data", "symbol", symbol, "limit", limit)

	var prices []domain.PriceData

	query := `
select open_time as date, close_price as price
from candles 
where symbol = upper(?)
order by close_time desc
limit ?
`

	err := p.db.Raw(query, symbol, limit).Scan(&prices).Error

	if err != nil {
		slog.Error("Failed to query latest candles data", "error", err)
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	if len(prices) == 0 {
		slog.Warn("No recent candles data found", "symbol", symbol)
		return nil, fmt.Errorf("no recent price data found for symbol: %s", symbol)
	}

	slog.Info("Latest daily price data retrieved", "symbol", symbol, "count", len(prices))
	return prices, nil
}
