package postgres

import (
	"github.com/daddydemir/crypto/pkg/model"
	"gorm.io/gorm"
)

type postgresExchangeRepository struct {
	db *gorm.DB
}

func NewPostgresExchangeRepository(db *gorm.DB) model.ExchangeRepository {
	return &postgresExchangeRepository{
		db: db,
	}
}

func (p *postgresExchangeRepository) FindAll() ([]model.ExchangeModel, error) {
	var exchanges []model.ExchangeModel
	p.db.Find(&exchanges)
	return exchanges, nil
}
