package app

import (
	"github.com/daddydemir/crypto/pkg/portfolio/domain"
	"github.com/daddydemir/crypto/pkg/portfolio/infra"
)

type App struct {
	repository *infra.Repository
}

func NewApp(repository *infra.Repository) *App {
	return &App{repository: repository}
}

func (a *App) List(username string) ([]domain.Transaction, error) {
	return a.repository.List(username)
}

func (a *App) Create(username string, transaction domain.Transaction) (*domain.Transaction, error) {
	transaction.Username = username
	if err := transaction.NormalizeAndValidate(); err != nil {
		return nil, err
	}
	if err := a.repository.Create(&transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (a *App) Update(username string, id uint, transaction domain.Transaction) (*domain.Transaction, error) {
	existing, err := a.repository.Find(username, id)
	if err != nil {
		return nil, err
	}
	existing.TransactionType = transaction.TransactionType
	existing.BaseAsset = transaction.BaseAsset
	existing.QuoteAsset = transaction.QuoteAsset
	existing.ReceivedAsset = transaction.ReceivedAsset
	existing.ReceivedAmount = transaction.ReceivedAmount
	existing.SpentAsset = transaction.SpentAsset
	existing.SpentAmount = transaction.SpentAmount
	existing.USDValue = transaction.USDValue
	existing.FeeAmount = transaction.FeeAmount
	existing.FeeAsset = transaction.FeeAsset
	existing.Platform = transaction.Platform
	existing.TradedAt = transaction.TradedAt
	existing.Notes = transaction.Notes
	if err = existing.NormalizeAndValidate(); err != nil {
		return nil, err
	}
	if err = a.repository.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (a *App) Import(username string, transactions []domain.Transaction) (int64, error) {
	for i := range transactions {
		transactions[i].Username = username
		if err := transactions[i].NormalizeAndValidate(); err != nil {
			return 0, err
		}
	}
	return a.repository.Import(transactions)
}

func (a *App) Delete(username string, id uint) error {
	return a.repository.Delete(username, id)
}
