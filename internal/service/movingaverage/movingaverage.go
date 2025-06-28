package movingaverage

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/service/chart"
	"log/slog"
	"time"

	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/signal"
	"github.com/google/uuid"
)

type Service struct {
	cache    port.Cache
	signal   signal.SignalService
	notifier port.Broker
	rsi      *chart.RsiService
}

func NewService(c port.Cache, s signal.SignalService, n port.Broker, r *chart.RsiService) *Service {
	return &Service{
		cache:    c,
		signal:   s,
		notifier: n,
		rsi:      r,
	}
}

func (s *Service) CheckAll(short, middle, long int) {
	coins := s.getCoins()
	var signals []model.SignalModel

	for index, coin := range coins {
		if index > 30 {
			continue
		}
		sItem, mItem, lItem, err := s.calculateItems(coin.Id, short, middle, long)
		if err != nil {
			slog.Error("calculateItems failed", "coin", coin.Id, "err", err)
			continue
		}
		if !sameDate(sItem, mItem, lItem) {
			slog.Warn("date mismatch", "coin", coin.Id)
			continue
		}
		pattern := evaluatePattern(sItem.Value, mItem.Value, lItem.Value)

		index := s.rsi.Index()

		signals = append(signals, model.SignalModel{
			ID:         uuid.New(),
			ExchangeId: coin.Id,
			CreateDate: time.Now(),
			Short:      sItem.Value,
			Middle:     mItem.Value,
			Long:       lItem.Value,
			Rsi:        index,
			Result:     pattern,
		})

		if pattern == "7 > 25 > 99" {
			s.notifier.SendMessage(fmt.Sprintf("[BUY] %s RSI: %.0f", coin.Id, index))
		} else if pattern == "99 > 25 > 7" {
			s.notifier.SendMessage(fmt.Sprintf("[SELL] %s RSI: %.0f", coin.Id, index))
		}
	}

	if err := s.signal.SaveAll(signals); err != nil {
		slog.Error("failed to save signals", "err", err)
	}
}

func (s *Service) calculateItems(coin string, short, middle, long int) (model.ChartModel, model.ChartModel, model.ChartModel, error) {

	sList := chart.NewEmaService(coin, short, nil).Calculate()
	mList := chart.NewEmaService(coin, middle, nil).Calculate()
	lList := chart.NewEmaService(coin, long, nil).Calculate()

	if len(sList) == 0 || len(mList) == 0 || len(lList) == 0 {
		return model.ChartModel{}, model.ChartModel{}, model.ChartModel{}, fmt.Errorf("insufficient data")
	}

	return model.ChartModel(sList[len(sList)-1]), model.ChartModel(mList[len(mList)-1]), model.ChartModel(lList[len(lList)-1]), nil
}

func (s *Service) getCoins() []model.Coin {
	data, err := s.cache.Get("coinList")
	if err != nil {
		slog.Error("getCoins: cache error", "err", err)
		return nil
	}
	var coins []model.Coin
	str, ok := data.(string)
	if !ok {
		slog.Error("getCoins: invalid cache format")
		return nil
	}
	if err = json.Unmarshal([]byte(str), &coins); err != nil {
		slog.Error("getCoins: unmarshal error", "err", err)
		return nil
	}
	return coins
}

func sameDate(a, b, c model.ChartModel) bool {
	return a.Date.Equal(b.Date) && b.Date.Equal(c.Date)
}

func evaluatePattern(short, middle, long float32) string {
	switch {
	case short > middle && middle > long:
		return "7 > 25 > 99"
	case long > short && short > middle:
		return "99 > 7 > 25"
	case middle > short && short > long:
		return "25 > 7 > 99"
	case middle > long && long > short:
		return "25 > 99 > 7"
	case short > long && long > middle:
		return "7 > 99 > 25"
	case long > middle && middle > short:
		return "99 > 25 > 7"
	default:
		return ""
	}
}
