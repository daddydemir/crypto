package infra

import (
	"context"
	"github.com/daddydemir/crypto/pkg/analyses/alert/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, a *domain.Alert) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repository) Update(ctx context.Context, a *domain.Alert) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Alert{}, id).Error
}

func (r *Repository) FindByID(ctx context.Context, id uint) (*domain.Alert, error) {
	var a domain.Alert
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Alert, error) {
	var alerts []domain.Alert
	if err := r.db.WithContext(ctx).Where("is_active = true").Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}
