package coin

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/domain/coin"
	"strings"
)

type GetTopCoinsStats struct {
	historyRepo coin.HistoryRepository
	marketRepo  coin.MarketRepository
}

func NewGetTopCoinsStats(historyRepo coin.HistoryRepository, marketRepo coin.MarketRepository) *GetTopCoinsStats {
	return &GetTopCoinsStats{
		historyRepo: historyRepo,
		marketRepo:  marketRepo,
	}
}

func (u *GetTopCoinsStats) Execute() ([]StatsDTO, error) {
	currentCoins, err := u.marketRepo.GetCurrentPrices()
	if err != nil {
		return nil, err
	}

	coinMap, err := convertToMap(u.marketRepo.GetPriceChanges())
	if err != nil {
		return nil, err
	}

	var results []StatsDTO

	for _, c := range currentCoins {
		results = append(results, StatsDTO{
			ID:        c.ID,
			Name:      c.Name,
			Symbol:    c.Symbol,
			Price:     c.PriceUSD,
			Change24h: coinMap[strings.ToLower(c.Symbol)].Change24h,
			Change7d:  coinMap[strings.ToLower(c.Symbol)].Change7d,
		})
	}

	return results, nil
}

func convertToMap(list []coin.PriceResult, err error) (map[string]coin.PriceResult, error) {
	if err != nil {
		return nil, fmt.Errorf("failed to convert to map: %w", err)
	}
	coinMap := make(map[string]coin.PriceResult)
	for _, c := range list {
		coinMap[c.ExchangeID] = c
	}

	return coinMap, err
}
