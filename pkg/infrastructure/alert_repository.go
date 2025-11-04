package infrastructure

import (
	"context"
	"github.com/daddydemir/crypto/pkg/domain/alert"
	"gorm.io/gorm"
)

type AlertRepositoryImpl struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepositoryImpl {
	return &AlertRepositoryImpl{db: db}
}

func (r *AlertRepositoryImpl) Save(ctx context.Context, a *alert.Alert) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AlertRepositoryImpl) Update(ctx context.Context, a *alert.Alert) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *AlertRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&alert.Alert{}, id).Error
}

func (r *AlertRepositoryImpl) FindByID(ctx context.Context, id uint) (*alert.Alert, error) {
	var a alert.Alert
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AlertRepositoryImpl) List(ctx context.Context) ([]alert.Alert, error) {
	var alerts []alert.Alert
	if err := r.db.WithContext(ctx).Where("is_active = true").Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}
