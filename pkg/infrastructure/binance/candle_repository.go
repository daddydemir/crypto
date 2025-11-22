package binance

import (
	"github.com/daddydemir/crypto/pkg/domain/binance"
	"gorm.io/gorm"
)

type CandleRepository struct {
	db *gorm.DB
}

func NewCandleRepository(db *gorm.DB) *CandleRepository {
	return &CandleRepository{db: db}
}

func (cr *CandleRepository) Save(candle binance.Candle) error {
	return cr.db.Create(&candle).Error
}

func (cr *CandleRepository) SaveMany(candles []binance.Candle) error {
	return cr.db.Create(&candles).Error
}
