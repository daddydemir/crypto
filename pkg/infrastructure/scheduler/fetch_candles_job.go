package scheduler

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/application/binance"
	"github.com/daddydemir/crypto/pkg/service"
	"github.com/robfig/cron/v3"
	"log/slog"
	"slices"
	"time"
)

var unnecessaryCoins = []string{
	"USDT", "USDC", "USDE", "USD1", "USDf", "USDG", "USDY",
	"BTCB", "WBTC", "EWTH", "WEETH", "STETH", "PYUSD", "WBNB",
	"RLUSD", "XAUT", "FDUSD",
}

func FetchCandlesJob(service *binance.CandleService, cacheService service.CacheService) *cron.Cron {
	location, _ := time.LoadLocation("Turkey")
	c := cron.New(cron.WithLocation(location))

	c.AddFunc("00 05 * * *", func() {
		yesterday := time.Now().Add(-1 * time.Hour * 24).UnixMilli()
		today := time.Now().UnixMilli()

		coins := cacheService.GetCoins()
		for _, coin := range coins {
			if !slices.Contains(unnecessaryCoins, coin.Symbol) {
				err := service.FetchAndStore(fmt.Sprintf("%sUSDT", coin.Symbol), "1d", yesterday, today, 1)
				if err != nil {
					slog.Error("FetchCandlesJob:service.FetchAndStore", "coin", coin.Symbol, "error", err)
				} else {
					slog.Info("FetchCandlesJob:service.FetchAndStore", "coin", coin.Symbol)
				}
			}
		}
	})
	return c
}
