package service

import (
	"github.com/daddydemir/crypto/pkg/model"
	"log/slog"
)

type ExchangeService struct {
	exchangeRepo model.ExchangeRepository
}

func NewExchangeService(exchangeRepo model.ExchangeRepository) *ExchangeService {
	return &ExchangeService{exchangeRepo: exchangeRepo}
}

func (e *ExchangeService) FindAll() ([]model.ExchangeModel, error) {
	return e.exchangeRepo.FindAll()
}

func (e *ExchangeService) SaveAll(exchanges []model.ExchangeModel) error {
	return e.exchangeRepo.SaveAll(exchanges)
}

func (e *ExchangeService) CreateExchange() {
	err := e.SaveAll(GetExchange())
	if err != nil {
		slog.Error("CreateExchange.SaveAll", "error", err)
	}
}
