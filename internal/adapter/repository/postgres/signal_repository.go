package postgres

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/signal"
	"gorm.io/gorm"
)

type PostgresSignalRepository struct {
	db *gorm.DB
}

func NewPostgresSignalRepository(db *gorm.DB) signal.SignalRepository {
	return &PostgresSignalRepository{db: db}
}

func (p *PostgresSignalRepository) SaveAll(signals []model.SignalModel) error {
	return p.db.Create(&signals).Error
}
