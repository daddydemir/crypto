package application

import "github.com/daddydemir/crypto/pkg/binance/domain"

type GetCandlesQuery struct {
	repo domain.Repository
}

func NewGetCandlesQuery(repo domain.Repository) *GetCandlesQuery {
	return &GetCandlesQuery{repo: repo}
}

func (q *GetCandlesQuery) Execute(symbol string) (candles []domain.Candle, err error) {
	return q.repo.GetBySymbol(symbol)
}
