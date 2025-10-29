package coin

import (
	"github.com/daddydemir/crypto/pkg/domain/coin"
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

	var results []StatsDTO

	for _, c := range currentCoins {
		// todo:
		//price24h, err := u.historyRepo.GetPriceAt(c.ID, 1)
		//if err != nil {
		//	continue
		//}
		//
		//price7d, err := u.historyRepo.GetPriceAt(c.ID, 7)
		//if err != nil {
		//	continue
		//}

		//fmt.Printf("coin: %s, price: %f, price24h: %f, price7d: %f\n", c.ID, c.PriceUSD, price24h, price7d)

		//change24h := ((c.PriceUSD - price24h) / price24h) * 100
		//change7d := ((c.PriceUSD - price7d) / price7d) * 100
		change24h := 0.0
		change7d := 0.0

		results = append(results, StatsDTO{
			ID:        c.ID,
			Name:      c.Name,
			Symbol:    c.Symbol,
			Price:     c.PriceUSD,
			Change24h: change24h,
			Change7d:  change7d,
		})
	}

	return results, nil
}
