package domain

import (
	"fmt"
	"strings"
	"time"
)

const (
	TransactionTypeBuy  = "BUY"
	TransactionTypeSell = "SELL"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Username        string    `gorm:"size:25;not null;index" json:"-"`
	CoinSymbol      string    `gorm:"size:20;not null;index" json:"coinSymbol"`
	TransactionType string    `gorm:"size:4;not null" json:"transactionType"`
	BaseAsset       string    `gorm:"size:20;not null;default:''" json:"baseAsset"`
	QuoteAsset      string    `gorm:"size:20;not null;default:''" json:"quoteAsset"`
	ReceivedAsset   string    `gorm:"size:20;not null;default:''" json:"receivedAsset"`
	ReceivedAmount  float64   `gorm:"type:numeric(30,10);not null;default:0" json:"receivedAmount"`
	SpentAsset      string    `gorm:"size:20;not null;default:''" json:"spentAsset"`
	SpentAmount     float64   `gorm:"type:numeric(30,10);not null;default:0" json:"spentAmount"`
	USDValue        float64   `gorm:"type:numeric(30,10)" json:"usdValue"`
	FeeAmount       float64   `gorm:"type:numeric(30,10)" json:"feeAmount"`
	FeeAsset        string    `gorm:"size:20" json:"feeAsset"`
	Source          string    `gorm:"size:30;not null;default:'MANUAL'" json:"source"`
	ExternalTradeID string    `gorm:"size:100" json:"externalTradeId"`
	ExchangeOrderID string    `gorm:"-" json:"-"`
	LastFillAt      time.Time `gorm:"-" json:"-"`
	Quantity        float64   `gorm:"type:numeric(30,10);not null;default:0" json:"-"`
	UnitPrice       float64   `gorm:"type:numeric(30,10);not null;default:0" json:"-"`
	Platform        string    `gorm:"size:80;not null" json:"platform"`
	TradedAt        time.Time `gorm:"not null;index" json:"tradedAt"`
	Notes           string    `gorm:"size:500" json:"notes"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (Transaction) TableName() string {
	return "portfolio_transactions"
}

func (t *Transaction) NormalizeAndValidate() error {
	t.CoinSymbol = strings.ToUpper(strings.TrimSpace(t.CoinSymbol))
	t.BaseAsset = strings.ToUpper(strings.TrimSpace(t.BaseAsset))
	t.QuoteAsset = strings.ToUpper(strings.TrimSpace(t.QuoteAsset))
	t.ReceivedAsset = strings.ToUpper(strings.TrimSpace(t.ReceivedAsset))
	t.SpentAsset = strings.ToUpper(strings.TrimSpace(t.SpentAsset))
	t.FeeAsset = strings.ToUpper(strings.TrimSpace(t.FeeAsset))
	t.Source = strings.ToUpper(strings.TrimSpace(t.Source))
	t.TransactionType = strings.ToUpper(strings.TrimSpace(t.TransactionType))
	t.Platform = strings.TrimSpace(t.Platform)
	t.Notes = strings.TrimSpace(t.Notes)
	if t.BaseAsset == "" || t.QuoteAsset == "" {
		return fmt.Errorf("base and quote assets are required")
	}
	if t.TransactionType != TransactionTypeBuy && t.TransactionType != TransactionTypeSell {
		return fmt.Errorf("transaction type must be BUY or SELL")
	}
	if t.ReceivedAsset == "" || t.ReceivedAmount <= 0 {
		return fmt.Errorf("received asset and amount are required")
	}
	if t.SpentAsset == "" || t.SpentAmount <= 0 {
		return fmt.Errorf("spent asset and amount are required")
	}
	if t.Platform == "" {
		return fmt.Errorf("platform is required")
	}
	if t.TradedAt.IsZero() {
		return fmt.Errorf("trade date is required")
	}
	t.CoinSymbol = t.BaseAsset
	if t.TransactionType == TransactionTypeBuy {
		t.Quantity = t.ReceivedAmount
	} else {
		t.Quantity = t.SpentAmount
	}
	if t.Quantity > 0 {
		t.UnitPrice = t.USDValue / t.Quantity
	}
	if t.Source == "" {
		t.Source = "MANUAL"
	}
	return nil
}
