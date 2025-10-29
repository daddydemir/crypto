package coin

import (
	"github.com/daddydemir/crypto/pkg/domain/coin"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
)

type MarketRepository struct {
	api *coincap.CachedClient
}

func NewCoinGeckoMarketRepository(api *coincap.CachedClient) *MarketRepository {
	return &MarketRepository{api: api}
}

func (r *MarketRepository) GetCurrentPrices() ([]coin.Coin, error) {
	err, coins := r.api.ListCoins()
	if err != nil {
		return nil, err
	}
	response := make([]coin.Coin, 0, 100)
	for _, c := range coins {
		response = append(response, coin.Coin{
			ID:       c.Id,
			Name:     c.Name,
			Symbol:   c.Symbol,
			PriceUSD: float64(c.PriceUsd),
		})
	}
	return response, nil
}
