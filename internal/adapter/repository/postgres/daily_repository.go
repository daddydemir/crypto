package postgres

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/daily"
	"gorm.io/gorm"
)

type PostgresDailyRepository struct {
	db *gorm.DB
}

func NewPostgresDailyRepository(db *gorm.DB) daily.DailyRepository {
	return &PostgresDailyRepository{db: db}
}

func (r *PostgresDailyRepository) FindByDateRange(start, end string) ([]model.DailyModel, error) {
	var dailies []model.DailyModel
	err := r.db.Where("date BETWEEN ? AND ?", start, end).
		Order("exchange_id, date").
		Find(&dailies).Error
	return dailies, err
}

func (r *PostgresDailyRepository) FindByIdAndDateRange(id, start, end string) ([]model.DailyModel, error) {
	var dailies []model.DailyModel
	err := r.db.Where("exchange_id = ? AND date BETWEEN ? AND ?", id, start, end).
		Order("date").
		Find(&dailies).Error
	return dailies, err
}

func (r *PostgresDailyRepository) FindTopSmallerByRate(start, end string) ([5]model.DailyModel, error) {
	var dailies [5]model.DailyModel
	err := r.db.Raw(`
		SELECT * FROM daily_models 
		WHERE date BETWEEN ? AND ? AND avg < 1 
		ORDER BY rate DESC LIMIT 5`, start, end).Scan(&dailies).Error
	return dailies, err
}

func (r *PostgresDailyRepository) FindTopBiggerByRate(start, end string) ([5]model.DailyModel, error) {
	var dailies [5]model.DailyModel
	err := r.db.Raw(`
		SELECT * FROM daily_models 
		WHERE date BETWEEN ? AND ? AND avg > 1 
		ORDER BY rate DESC LIMIT 5`, start, end).Scan(&dailies).Error
	return dailies, err
}

func (r *PostgresDailyRepository) SaveAll(models []model.DailyModel) error {
	return r.db.Save(&models).Error
}
