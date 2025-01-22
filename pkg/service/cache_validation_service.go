package service

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
	"time"
)

type ValidateService struct {
	cache.Cache
	client coincap.CoinCapClient
}

func NewValidateService(c cache.Cache, client coincap.CoinCapClient) *ValidateService {
	return &ValidateService{
		c,
		client,
	}
}

func (v *ValidateService) Validate() {
	cacheService := v.Cache
	client := v.client
	var coins []coincap.Coin
	data, err := cacheService.Get("coinList")
	if err != nil {
		slog.Error("Validate:cacheService.Get", "key", "coinList", "err", err)
		return
	}

	bytes, ok := data.(string)
	if !ok {
		slog.Error("Validate:data.(string)", "data", data)
		return
	}

	err = json.Unmarshal([]byte(bytes), &coins)
	if err != nil {
		slog.Error("Validate:json.Unmarshal", "bytes", bytes, "err", err)
		return
	}

	for _, i := range coins {
		array := make([]coincap.History, 0, 1)
		err = cacheService.GetList(i.Id, &array, -1, -1)
		if err != nil {
			slog.Error("Validate:cacheService.GetList", "coin", i.Id, "error", err)
			continue
		}
		_, histories := client.HistoryWithTime(i.Id, array[0].Date.Add(time.Hour*24).UnixNano(), time.Now().UnixNano())

		if len(histories) == 0 {
			slog.Error("Validate:coincap.HistoryWithTime", "coin", i.Id, "error", "list is empty")
			continue
		}

		err = cacheService.SetList(i.Id, histories, 0)
		if err != nil {
			slog.Error("Validate:cacheService.SetList", "error", err)
			continue
		}
		slog.Info("Validation Success", "coin", i.Id)
		time.Sleep(time.Second * 20)
	}
}
