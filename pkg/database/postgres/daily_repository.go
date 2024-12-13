package postgres

import (
	"github.com/daddydemir/crypto/pkg/model"
	"gorm.io/gorm"
)

type postgresDailyRepository struct {
	db *gorm.DB
}

func NewPostgresDailyRepository(db *gorm.DB) model.DailyRepository {
	return &postgresDailyRepository{
		db: db,
	}
}

func (r *postgresDailyRepository) FindByDateRange(start, end string) ([]model.DailyModel, error) {
	var dailies []model.DailyModel
	tx := r.db.Where("date between ?  and ? order by exchange_id , date", start, end).Find(&dailies)
	return dailies, tx.Error
}

func (r *postgresDailyRepository) FindByIdAndDateRange(id, start, end string) ([]model.DailyModel, error) {
	var dailies []model.DailyModel
	tx := r.db.Where("date between ? and ? and exchange_id = ?", start, end, id).Find(&dailies)
	return dailies, tx.Error
}

func (r *postgresDailyRepository) FindTopSmallerByRate(start, end string) ([5]model.DailyModel, error) {
	var dailies [5]model.DailyModel
	tx := r.db.Where("date between ? and ? and avg < 1 order by rate desc limit 5", start, end).Find(&dailies)
	return dailies, tx.Error
}

func (r *postgresDailyRepository) FindTopBiggerByRate(start, end string) ([5]model.DailyModel, error) {
	var dailies [5]model.DailyModel
	tx := r.db.Where("date between ? and ? and avg > 1 order by rate desc limit 5", start, end).Find(&dailies)
	return dailies, tx.Error
}

func (r *postgresDailyRepository) SaveAll(models []model.DailyModel) error {
	tx := r.db.Save(&models)
	return tx.Error
}
