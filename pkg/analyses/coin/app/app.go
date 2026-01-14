package app

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/analyses/coin/domain"
	"github.com/daddydemir/crypto/pkg/analyses/coin/infra"
	"strings"
)

type App struct {
	repo *infra.Repository
}

func NewApp(repo *infra.Repository) *App {
	return &App{repo: repo}
}

func (a *App) GetTopCoins() ([]domain.StatsDTO, error) {
	currentCoins, err := a.repo.GetCurrentPrices()
	if err != nil {
		return nil, err
	}

	coinMap, err := convertToMap(a.repo.GetPriceChanges())
	if err != nil {
		return nil, err
	}

	var results []domain.StatsDTO

	for _, c := range currentCoins {
		results = append(results, domain.StatsDTO{
			ID:                  c.ID,
			Name:                c.Name,
			Symbol:              c.Symbol,
			Price:               c.PriceUSD,
			Change24h:           coinMap[strings.ToLower(c.Symbol)].Change24h,
			Change7d:            coinMap[strings.ToLower(c.Symbol)].Change7d,
			Change30d:           coinMap[strings.ToLower(c.Symbol)].Change30d,
			ArithmeticChange7d:  coinMap[strings.ToLower(c.Symbol)].ChangeArithmetic7d,
			ArithmeticChange30d: coinMap[strings.ToLower(c.Symbol)].ChangeArithmetic30d,
		})
	}

	return results, nil
}

func convertToMap(list []domain.PriceResult, err error) (map[string]domain.PriceResult, error) {
	if err != nil {
		return nil, fmt.Errorf("failed to convert to map: %w", err)
	}
	coinMap := make(map[string]domain.PriceResult)
	for _, c := range list {
		coinMap[c.ExchangeID] = c
	}

	return coinMap, err
}
