package port

import "github.com/daddydemir/crypto/internal/domain/model"

type CoingeckoClient interface {
	GetTopHundred() ([]model.ExchangeModel, error)
	GetTopHundredDaily() ([]model.DailyModel, error)
}
