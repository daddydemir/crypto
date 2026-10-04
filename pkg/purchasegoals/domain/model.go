package domain

import (
	"fmt"
	"strings"
	"time"
)

const (
	StatusWanted    = "WANTED"
	StatusSaving    = "SAVING"
	StatusPurchased = "PURCHASED"
	StatusCancelled = "CANCELLED"
)

type Goal struct {
	ID                    uint                `gorm:"primaryKey" json:"id"`
	Username              string              `gorm:"size:25;not null;index" json:"-"`
	Name                  string              `gorm:"size:160;not null" json:"name"`
	Notes                 string              `gorm:"size:2000" json:"notes"`
	Currency              string              `gorm:"size:3;not null" json:"currency"`
	Priority              string              `gorm:"size:20;not null" json:"priority"`
	Status                string              `gorm:"size:20;not null;index" json:"status"`
	TargetDate            *time.Time          `json:"targetDate"`
	ExpectedMonthlySaving float64             `gorm:"type:numeric(20,2);not null;default:0" json:"expectedMonthlySaving"`
	Items                 []Item              `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"items"`
	Savings               []SavingTransaction `gorm:"foreignKey:GoalID;constraint:OnDelete:CASCADE" json:"savings"`
	CreatedAt             time.Time           `json:"createdAt"`
	UpdatedAt             time.Time           `json:"updatedAt"`
}

func (Goal) TableName() string { return "purchase_goals" }

type Item struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	GoalID              uint           `gorm:"not null;index" json:"goalId"`
	Name                string         `gorm:"size:160;not null" json:"name"`
	Description         string         `gorm:"size:2000" json:"description"`
	CurrentPrice        float64        `gorm:"type:numeric(20,2);not null" json:"currentPrice"`
	TargetPrice         *float64       `gorm:"type:numeric(20,2)" json:"targetPrice"`
	Currency            string         `gorm:"size:3;not null" json:"currency"`
	ImageURL            string         `gorm:"size:2000" json:"imageUrl"`
	ProductURL          string         `gorm:"size:2000" json:"productUrl"`
	Category            string         `gorm:"size:100" json:"category"`
	Priority            string         `gorm:"size:20;not null" json:"priority"`
	Status              string         `gorm:"size:20;not null;index" json:"status"`
	TargetDate          *time.Time     `json:"targetDate"`
	ImageFit            string         `gorm:"size:10;not null;default:'COVER'" json:"imageFit"`
	ImagePositionX      int            `gorm:"not null;default:50" json:"imagePositionX"`
	ImagePositionY      int            `gorm:"not null;default:50" json:"imagePositionY"`
	ImageZoom           float64        `gorm:"type:numeric(4,2);not null;default:1" json:"imageZoom"`
	ActualPurchasePrice *float64       `gorm:"type:numeric(20,2)" json:"actualPurchasePrice"`
	PurchasedAt         *time.Time     `json:"purchasedAt"`
	PriceHistory        []PriceHistory `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE" json:"priceHistory"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
}

func (Item) TableName() string { return "purchase_items" }

type SavingTransaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GoalID    uint      `gorm:"not null;index" json:"goalId"`
	Amount    float64   `gorm:"type:numeric(20,2);not null" json:"amount"`
	Date      time.Time `gorm:"not null;index" json:"date"`
	Note      string    `gorm:"size:500" json:"note"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (SavingTransaction) TableName() string { return "purchase_savings" }

type PriceHistory struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	ItemID   uint      `gorm:"not null;index" json:"itemId"`
	Price    float64   `gorm:"type:numeric(20,2);not null" json:"price"`
	Currency string    `gorm:"size:3;not null" json:"currency"`
	Date     time.Time `gorm:"not null;index" json:"date"`
}

func (PriceHistory) TableName() string { return "purchase_price_history" }

type GoalView struct {
	Goal
	TargetAmount       float64  `json:"targetAmount"`
	SavedAmount        float64  `json:"savedAmount"`
	RemainingAmount    float64  `json:"remainingAmount"`
	Progress           float64  `json:"progress"`
	DisplayStatus      string   `json:"displayStatus"`
	PurchasedValue     float64  `json:"purchasedValue"`
	RemainingItemValue float64  `json:"remainingItemValue"`
	EstimatedMonths    *float64 `json:"estimatedMonths"`
}

func (g *Goal) NormalizeAndValidate() error {
	g.Name, g.Notes = strings.TrimSpace(g.Name), strings.TrimSpace(g.Notes)
	g.Currency, g.Priority, g.Status = strings.ToUpper(strings.TrimSpace(g.Currency)), strings.ToUpper(strings.TrimSpace(g.Priority)), strings.ToUpper(strings.TrimSpace(g.Status))
	if g.Name == "" {
		return fmt.Errorf("name is required")
	}
	if g.Currency == "" {
		g.Currency = "TRY"
	}
	if g.Priority == "" {
		g.Priority = "MEDIUM"
	}
	if g.Status == "" {
		g.Status = StatusWanted
	}
	if g.ExpectedMonthlySaving < 0 {
		return fmt.Errorf("expected monthly saving cannot be negative")
	}
	return nil
}

func (i *Item) NormalizeAndValidate(currency string) error {
	i.Name, i.Description, i.ImageURL, i.ProductURL, i.Category = strings.TrimSpace(i.Name), strings.TrimSpace(i.Description), strings.TrimSpace(i.ImageURL), strings.TrimSpace(i.ProductURL), strings.TrimSpace(i.Category)
	i.Currency, i.Priority, i.Status, i.ImageFit = strings.ToUpper(strings.TrimSpace(i.Currency)), strings.ToUpper(strings.TrimSpace(i.Priority)), strings.ToUpper(strings.TrimSpace(i.Status)), strings.ToUpper(strings.TrimSpace(i.ImageFit))
	if i.Name == "" {
		return fmt.Errorf("name is required")
	}
	if i.CurrentPrice < 0 || (i.TargetPrice != nil && *i.TargetPrice < 0) {
		return fmt.Errorf("prices cannot be negative")
	}
	if i.Currency == "" {
		i.Currency = currency
	}
	if i.Currency != currency {
		return fmt.Errorf("item currency must match goal currency")
	}
	if i.Priority == "" {
		i.Priority = "MEDIUM"
	}
	if i.Status == "" {
		i.Status = StatusWanted
	}
	if i.ImageFit == "" {
		i.ImageFit = "COVER"
	}
	if i.ImageFit != "COVER" && i.ImageFit != "CONTAIN" {
		return fmt.Errorf("image fit must be COVER or CONTAIN")
	}
	if i.ImagePositionX < 0 || i.ImagePositionX > 100 || i.ImagePositionY < 0 || i.ImagePositionY > 100 {
		return fmt.Errorf("image position must be between 0 and 100")
	}
	if i.ImageZoom == 0 {
		i.ImageZoom = 1
	}
	if i.ImageZoom < 1 || i.ImageZoom > 3 {
		return fmt.Errorf("image zoom must be between 1 and 3")
	}
	return nil
}

func ItemTarget(i Item) float64 {
	if i.TargetPrice != nil {
		return *i.TargetPrice
	}
	return i.CurrentPrice
}
