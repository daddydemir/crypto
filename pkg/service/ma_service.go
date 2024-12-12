package service

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/ma"
	"github.com/daddydemir/crypto/pkg/graphs/rsi"
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type MovingAverageService struct {
	cache         cache.Cache
	broker        broker.Broker
	signalService SignalService
}

func NewMovingAverageService(c cache.Cache, b broker.Broker, s SignalService) *MovingAverageService {
	return &MovingAverageService{
		cache:         c,
		broker:        b,
		signalService: s,
	}
}

func (m *MovingAverageService) CheckWithId(coin string, short, middle, long int) {

	shortService := ma.NewEma(coin, short)
	middleService := ma.NewEma(coin, middle)
	longService := ma.NewEma(coin, long)

	shortList := shortService.Calculate()
	middleList := middleService.Calculate()
	longList := longService.Calculate()

	shortItem := shortList[len(shortList)-1]
	middleItem := middleList[len(middleList)-1]
	longItem := longList[len(longList)-1]

	length := len(shortList) - 1

	if shortItem.Date == longItem.Date && shortItem.Date == middleItem.Date {

		for i := 0; i < 200; i++ {
			shortItem = shortList[length-i]
			middleItem = middleList[length-i]
			longItem = longList[length-i]
			fmt.Printf("date: %-10v : short: %-15v, middle: %-15v, long: %-15v, result: %-15v \n",
				shortItem.Date.Format("2006-01-02"),
				shortItem.Value,
				middleItem.Value,
				longItem.Value,
				result(shortItem.Value, middleItem.Value, longItem.Value),
			)

		}
	} else {
		println("Veri hatasi mevcut...")
	}

}

func (m *MovingAverageService) CheckAll(short, middle, long int) {
	signals := make([]model.SignalModel, 0)
	coins := m.getCoins()

	for index, coin := range coins {
		if index > 30 {
			continue
		}
		shortItem, middleItem, longItem, err := m.calculateMaItems(coin.Id, short, middle, long)
		if err != nil {
			slog.Error("CheckAll:calculateMaItems", "coin", coin.Id, "error", err)
			continue
		}

		if shortItem.Date == middleItem.Date && middleItem.Date == longItem.Date {
			r := result(shortItem.Value, middleItem.Value, longItem.Value)
			rsiService := rsi.NewRsi(coin.Id)
			rsiIndex := rsiService.Index()
			signals = append(signals, m.createSignalModel(coin.Id, r, shortItem.Value, middleItem.Value, longItem.Value, rsiIndex))
			if r == "7 > 25 > 99" {
				message := fmt.Sprintf("[Moving Average] \ncoin: %v Rsi: %0.f \nAsiri Alis", coin.Id, rsiIndex)
				slog.Info("CheckAll", "message", message)
				m.broker.SendMessage(message)
			} else if r == "99 > 25 > 7" {
				message := fmt.Sprintf("[Moving Average] \ncoin: %v Rsi: %0.f \nAsiri Satis", coin.Id, rsiIndex)
				slog.Info("CheckAll", "message", message)
				m.broker.SendMessage(message)
			} else {
				slog.Info("CheckAll", "message", "Onemsiz", "coin", coin.Id, "pattern", r)
			}
		} else {
			slog.Error("CheckAll", "coin", coin.Id, "message", "invalid data")
		}
	}
	err := m.signalService.SaveAll(signals)
	if err != nil {
		slog.Error("CheckAll:signalService.SaveAll", "error", err)
	}
}

func result(short, middle, long float32) string {

	if short > middle && middle > long {
		return "7 > 25 > 99"
	}

	if long > short && short > middle {
		return "99 > 7 > 25"
	}

	if middle > short && short > long {
		return "25 > 7 > 99"
	}

	if middle > long && long > short {
		return "25 > 99 > 7"
	}

	if short > long && long > middle {
		return "7 > 99 > 25"
	}

	if long > middle && middle > short {
		return "99 > 25 > 7"
	}

	return ""
}

func (m *MovingAverageService) getCoins() []coincap.Coin {
	cacheKey := "coinList"
	var coins []coincap.Coin
	cacheBody, err := m.cache.Get(cacheKey)
	if err != nil {
		slog.Error("CheckAll:Cache.Get", "key", cacheKey, "error", err)
		return nil
	}

	bytes, ok := cacheBody.(string)
	if !ok {
		slog.Error("Validate:data.(string)", "data", cacheBody)
		return nil
	}

	err = json.Unmarshal([]byte(bytes), &coins)
	if err != nil {
		slog.Error("Validate:json.Unmarshal", "bytes", bytes, "err", err)
		return nil
	}
	return coins
}

func (m *MovingAverageService) calculateMaItems(coin string, short, middle, long int) (shortItem, middleItem, longItem graphs.ChartModel, err error) {
	shortService := ma.NewEma(coin, short)
	middleService := ma.NewEma(coin, middle)
	longService := ma.NewEma(coin, long)

	shortList := shortService.Calculate()
	middleList := middleService.Calculate()
	longList := longService.Calculate()

	if len(shortList) == 0 || len(middleList) == 0 || len(longList) == 0 {
		err = fmt.Errorf("insufficient data for MA calculation")
		return
	}

	shortItem = shortList[len(shortList)-1]
	middleItem = middleList[len(middleList)-1]
	longItem = longList[len(longList)-1]
	return
}

func (m *MovingAverageService) createSignalModel(coin, result string, short, midd, long, rsi float32) model.SignalModel {
	return model.SignalModel{
		ID:         uuid.New(),
		ExchangeId: coin,
		CreateDate: time.Now(),
		Short:      short,
		Middle:     midd,
		Long:       long,
		Rsi:        rsi,
		Result:     result,
	}
}
