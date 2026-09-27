package infra

import (
	"errors"

	"github.com/daddydemir/crypto/pkg/portfolio/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	if err := r.db.AutoMigrate(&domain.Transaction{}); err != nil {
		return err
	}
	if err := r.db.Exec(`update portfolio_transactions set
		base_asset = coin_symbol, quote_asset = 'USD',
		received_asset = case when transaction_type = 'BUY' then coin_symbol else 'USD' end,
		received_amount = case when transaction_type = 'BUY' then quantity else quantity * unit_price end,
		spent_asset = case when transaction_type = 'BUY' then 'USD' else coin_symbol end,
		spent_amount = case when transaction_type = 'BUY' then quantity * unit_price else quantity end,
		usd_value = quantity * unit_price
		where base_asset = ''`).Error; err != nil {
		return err
	}
	return r.db.Exec(`create unique index if not exists ux_portfolio_external_trade
		on portfolio_transactions (username, source, external_trade_id, base_asset, quote_asset)
		where external_trade_id <> ''`).Error
}

func (r *Repository) List(username string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Where("username = ?", username).Order("traded_at desc, id desc").Find(&transactions).Error
	return transactions, err
}

func (r *Repository) Create(transaction *domain.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *Repository) Find(username string, id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.Where("id = ? and username = ?", id, username).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *Repository) Update(transaction *domain.Transaction) error {
	result := r.db.Model(&domain.Transaction{}).
		Where("id = ? and username = ?", transaction.ID, transaction.Username).
		Updates(map[string]any{
			"coin_symbol": transaction.CoinSymbol, "transaction_type": transaction.TransactionType,
			"quantity": transaction.Quantity, "unit_price": transaction.UnitPrice,
			"base_asset": transaction.BaseAsset, "quote_asset": transaction.QuoteAsset,
			"received_asset": transaction.ReceivedAsset, "received_amount": transaction.ReceivedAmount,
			"spent_asset": transaction.SpentAsset, "spent_amount": transaction.SpentAmount,
			"usd_value": transaction.USDValue, "fee_amount": transaction.FeeAmount, "fee_asset": transaction.FeeAsset,
			"platform": transaction.Platform, "traded_at": transaction.TradedAt, "notes": transaction.Notes,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) Import(transactions []domain.Transaction) (int64, error) {
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&transactions)
	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(username string, id uint) error {
	result := r.db.Where("id = ? and username = ?", id, username).Delete(&domain.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("portfolio transaction not found")
	}
	return nil
}
