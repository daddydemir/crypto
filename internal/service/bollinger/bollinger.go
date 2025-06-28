package bollinger

import (
	"fmt"
	"github.com/daddydemir/crypto/internal/port"
	"log/slog"

	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/chart"
)

type Service struct {
	client   port.CoinCapAPI
	notifier port.Broker
	builder  func(symbol string, period int) chart.BollingerChart
}

func NewService(
	client port.CoinCapAPI,
	notifier port.Broker,
	builder func(string, int) chart.BollingerChart,
) *Service {
	return &Service{
		client:   client,
		notifier: notifier,
		builder:  builder,
	}
}

func (s *Service) CheckThresholds() {
	err, coins := s.client.ListCoins()
	if err != nil || len(coins) == 0 {
		slog.Error("CheckThresholds:ListCoins", "error", err)
		return
	}

	for i, coin := range coins {
		if i > 30 {
			break
		}

		calculator := s.builder(coin.Id, 20)
		middle, upper, lower, original := calculator.CalculateBands()

		if len(upper) == 0 || len(lower) == 0 || len(middle) < 2 || len(original) == 0 {
			slog.Warn("Insufficient data", "coin", coin.Id)
			continue
		}

		upperVal := upper[len(upper)-1].Value
		lowerVal := lower[len(lower)-1].Value
		latestOriginal := original[len(original)-1].Value

		if coin.PriceUsd > upperVal {
			msg := fmt.Sprintf("[Bollinger ↑] %s > upper: %.3f > %.3f", coin.Id, coin.PriceUsd, upperVal)
			_ = s.notifier.SendMessage(msg)
		} else if coin.PriceUsd < lowerVal {
			msg := fmt.Sprintf("[Bollinger ↓] %s < lower: %.3f < %.3f", coin.Id, coin.PriceUsd, lowerVal)
			_ = s.notifier.SendMessage(msg)
		}

		if isCrossingDown(latestOriginal, coin.PriceUsd, middle) {
			msg := fmt.Sprintf("[Bollinger Trend ↓] %s düşüşte", coin.Id)
			_ = s.notifier.SendMessage(msg)
		}
	}
}

func isCrossingDown(prevPrice, currentPrice float32, middle []model.ChartModel) bool {
	prevMiddle := middle[len(middle)-2].Value
	currentMiddle := middle[len(middle)-1].Value
	return currentMiddle > currentPrice && prevPrice > prevMiddle
}
