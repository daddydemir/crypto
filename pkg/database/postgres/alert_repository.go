package postgres

import (
	"github.com/daddydemir/crypto/pkg/model"
	"gorm.io/gorm"
)

type postgresAlertRepository struct {
	db *gorm.DB
}

func NewPostgresAlertRepository(db *gorm.DB) model.AlertRepository {
	return &postgresAlertRepository{
		db: db,
	}
}

func (p *postgresAlertRepository) Save(alert model.Alert) error {
	tx := p.db.Create(&alert)
	return tx.Error
}

func (p *postgresAlertRepository) GetAll() ([]model.Alert, error) {
	var list []model.Alert
	tx := p.db.Where("is_active = true").Find(&list)
	return list, tx.Error
}

func (p *postgresAlertRepository) Delete(id int) error {
	tx := p.db.Where("id = ?", id).Delete(&model.Alert{})
	return tx.Error
}
