package application

import "github.com/daddydemir/crypto/pkg/binance/domain"

type GetCandlesQuery struct {
	repo domain.Repository
}

func NewGetCandlesQuery(repo domain.Repository) *GetCandlesQuery {
	return &GetCandlesQuery{repo: repo}
}

func (q *GetCandlesQuery) Execute(symbol, year, month string) (candles []domain.Candle, err error) {

	if year != "" && month != "" {
		return q.repo.GetBySymbolAndYearMonth(symbol, year, month)
	} else if year != "" {
		return q.repo.GetBySymbolAndYear(symbol, year)
	}
	return q.repo.GetBySymbol(symbol)
}
