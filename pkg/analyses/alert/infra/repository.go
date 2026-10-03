package infra

import (
	"context"
	"github.com/daddydemir/crypto/pkg/analyses/alert/domain"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migrate() error {
	if err := r.db.AutoMigrate(&domain.Alert{}); err != nil {
		return err
	}
	// Alert evaluation did not exist before ownership was added. Keep legacy rows
	// out of the worker and undo any accidental legacy trigger during rollout.
	return r.db.Model(&domain.Alert{}).Where("username = '' AND triggered_at IS NOT NULL").Updates(map[string]any{"is_active": true, "triggered_at": nil, "triggered_price": 0}).Error
}

func (r *Repository) Save(ctx context.Context, a *domain.Alert) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repository) Update(ctx context.Context, a *domain.Alert) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *Repository) Deactivate(ctx context.Context, username string, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.Alert{}).Where("id = ? and username = ?", id, username).Update("is_active", false).Error
}

func (r *Repository) FindByID(ctx context.Context, username string, id uint) (*domain.Alert, error) {
	var a domain.Alert
	if err := r.db.WithContext(ctx).Where("id = ? and username = ?", id, username).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) List(ctx context.Context, username string) ([]domain.Alert, error) {
	var alerts []domain.Alert
	if err := r.db.WithContext(ctx).Where("username = ?", username).Order("create_date desc").Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *Repository) SetStatus(ctx context.Context, username string, id uint, active bool) error {
	updates := map[string]any{"is_active": active}
	if active {
		updates["triggered_at"] = nil
		updates["triggered_price"] = 0
	}
	return r.db.WithContext(ctx).Model(&domain.Alert{}).Where("id = ? and username = ?", id, username).Updates(updates).Error
}

func (r *Repository) Active(ctx context.Context) ([]domain.Alert, error) {
	var alerts []domain.Alert
	err := r.db.WithContext(ctx).Where("is_active = true and username <> ''").Find(&alerts).Error
	return alerts, err
}

func (r *Repository) MarkTriggered(ctx context.Context, id uint, price float32, at time.Time) (bool, error) {
	result := r.db.WithContext(ctx).Model(&domain.Alert{}).Where("id = ? and is_active = true and username <> ''", id).Updates(map[string]any{"is_active": false, "triggered_at": at, "triggered_price": price})
	return result.RowsAffected == 1, result.Error
}
