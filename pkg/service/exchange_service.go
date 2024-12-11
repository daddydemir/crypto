package service

import "github.com/daddydemir/crypto/pkg/model"

type ExchangeService struct {
	exchangeRepo model.ExchangeRepository
}

func NewExchangeService(exchangeRepo model.ExchangeRepository) *ExchangeService {
	return &ExchangeService{exchangeRepo: exchangeRepo}
}

func (e *ExchangeService) FindAll() ([]model.ExchangeModel, error) {
	return e.exchangeRepo.FindAll()
}
