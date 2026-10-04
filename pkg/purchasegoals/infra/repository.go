package infra

import (
	"errors"
	"github.com/daddydemir/crypto/pkg/purchasegoals/domain"
	"gorm.io/gorm"
	"time"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(&domain.Goal{}, &domain.Item{}, &domain.SavingTransaction{}, &domain.PriceHistory{})
}
func (r *Repository) List(user string) ([]domain.Goal, error) {
	var v []domain.Goal
	e := r.db.Where("username = ?", user).Preload("Items.PriceHistory", func(db *gorm.DB) *gorm.DB { return db.Order("date asc") }).Preload("Savings", func(db *gorm.DB) *gorm.DB { return db.Order("date desc, id desc") }).Order("created_at desc").Find(&v).Error
	return v, e
}
func (r *Repository) Find(user string, id uint) (*domain.Goal, error) {
	var v domain.Goal
	e := r.db.Where("id=? and username=?", id, user).Preload("Items.PriceHistory", func(db *gorm.DB) *gorm.DB { return db.Order("date asc") }).Preload("Savings", func(db *gorm.DB) *gorm.DB { return db.Order("date desc, id desc") }).First(&v).Error
	return &v, e
}
func (r *Repository) Create(g *domain.Goal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		items := g.Items
		g.Items = nil
		g.Savings = nil
		if e := tx.Create(g).Error; e != nil {
			return e
		}
		for n := range items {
			items[n].ID = 0
			items[n].PriceHistory = nil
			items[n].GoalID = g.ID
			if e := tx.Create(&items[n]).Error; e != nil {
				return e
			}
			if e := tx.Create(&domain.PriceHistory{ItemID: items[n].ID, Price: items[n].CurrentPrice, Currency: items[n].Currency, Date: time.Now()}).Error; e != nil {
				return e
			}
		}
		g.Items = items
		return nil
	})
}
func (r *Repository) SaveGoal(g *domain.Goal) error {
	return r.db.Model(&domain.Goal{}).Where("id=? and username=?", g.ID, g.Username).Updates(map[string]any{"name": g.Name, "notes": g.Notes, "currency": g.Currency, "priority": g.Priority, "status": g.Status, "target_date": g.TargetDate, "expected_monthly_saving": g.ExpectedMonthlySaving}).Error
}
func (r *Repository) DeleteGoal(user string, id uint) error {
	res := r.db.Where("id=? and username=?", id, user).Delete(&domain.Goal{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *Repository) FindItem(user string, id uint) (*domain.Item, *domain.Goal, error) {
	var item domain.Item
	e := r.db.Joins("join purchase_goals on purchase_goals.id = purchase_items.goal_id").Where("purchase_items.id=? and purchase_goals.username=?", id, user).First(&item).Error
	if e != nil {
		return nil, nil, e
	}
	var g domain.Goal
	e = r.db.First(&g, item.GoalID).Error
	return &item, &g, e
}
func (r *Repository) CreateItem(i *domain.Item) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Create(i).Error; e != nil {
			return e
		}
		return tx.Create(&domain.PriceHistory{ItemID: i.ID, Price: i.CurrentPrice, Currency: i.Currency, Date: time.Now()}).Error
	})
}
func (r *Repository) SaveItem(i *domain.Item) error { return r.db.Save(i).Error }
func (r *Repository) DeleteItem(user string, id uint) error {
	item, _, e := r.FindItem(user, id)
	if e != nil {
		return e
	}
	return r.db.Delete(item).Error
}
func (r *Repository) UpdatePrice(i *domain.Item, price float64, currency string, date time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Model(i).Update("current_price", price).Error; e != nil {
			return e
		}
		return tx.Create(&domain.PriceHistory{ItemID: i.ID, Price: price, Currency: currency, Date: date}).Error
	})
}
func (r *Repository) AddSaving(s *domain.SavingTransaction) error { return r.db.Create(s).Error }
func (r *Repository) FindSaving(user string, goalID, id uint) (*domain.SavingTransaction, error) {
	var s domain.SavingTransaction
	e := r.db.Joins("join purchase_goals on purchase_goals.id=purchase_savings.goal_id").Where("purchase_savings.id=? and purchase_savings.goal_id=? and purchase_goals.username=?", id, goalID, user).First(&s).Error
	return &s, e
}
func (r *Repository) SaveSaving(s *domain.SavingTransaction) error { return r.db.Save(s).Error }
func (r *Repository) DeleteSaving(user string, goalID, id uint) error {
	s, e := r.FindSaving(user, goalID, id)
	if e != nil {
		return e
	}
	return r.db.Delete(s).Error
}
func (r *Repository) SavingTotal(goalID uint, excludeID uint) (float64, error) {
	var total float64
	q := r.db.Model(&domain.SavingTransaction{}).Where("goal_id=?", goalID)
	if excludeID > 0 {
		q = q.Where("id<>?", excludeID)
	}
	e := q.Select("coalesce(sum(amount),0)").Scan(&total).Error
	return total, e
}

var ErrNegativeSavings = errors.New("total saved amount cannot be negative")
