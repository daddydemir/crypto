package postgres

import (
	"github.com/daddydemir/crypto/pkg/model"
	"gorm.io/gorm"
)

type postgresSignalRepository struct {
	db *gorm.DB
}

func NewPostgresSignalRepository(db *gorm.DB) model.SignalRepository {
	return &postgresSignalRepository{db: db}
}

func (p postgresSignalRepository) SaveAll(signals []model.SignalModel) error {
	p.db.Save(&signals)
	return nil
}
