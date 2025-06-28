package exchange

import "github.com/daddydemir/crypto/internal/domain/model"

type ExchangeRepository interface {
	FindAll() ([]model.ExchangeModel, error)
	SaveAll([]model.ExchangeModel) error
}
