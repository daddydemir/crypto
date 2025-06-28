package postgres

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/alert"
	"gorm.io/gorm"
)

type PostgresAlertRepository struct {
	db *gorm.DB
}

func NewPostgresAlertRepository(db *gorm.DB) alert.AlertRepository {
	return &PostgresAlertRepository{db: db}
}

func (p *PostgresAlertRepository) Save(alert model.Alert) error {
	return p.db.Create(&alert).Error
}

func (p *PostgresAlertRepository) GetAll() ([]model.Alert, error) {
	var alerts []model.Alert
	err := p.db.Where("is_active = true").Find(&alerts).Error
	return alerts, err
}

func (p *PostgresAlertRepository) Delete(id int) error {
	return p.db.Where("id = ?", id).Delete(&model.Alert{}).Error
}
