package postgres

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/exchange"
	"gorm.io/gorm"
)

type PostgresExchangeRepository struct {
	db *gorm.DB
}

func NewPostgresExchangeRepository(db *gorm.DB) exchange.ExchangeRepository {
	return &PostgresExchangeRepository{db: db}
}

func (p *PostgresExchangeRepository) FindAll() ([]model.ExchangeModel, error) {
	var exchanges []model.ExchangeModel
	err := p.db.Find(&exchanges).Error
	return exchanges, err
}

func (p *PostgresExchangeRepository) SaveAll(models []model.ExchangeModel) error {
	return p.db.Save(&models).Error
}
