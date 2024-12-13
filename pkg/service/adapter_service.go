package service

import (
	"log/slog"

	"github.com/daddydemir/crypto/pkg/adapter"
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/daddydemir/crypto/pkg/remote/coingecko"
)

func GetExchange() []model.ExchangeModel {
	var adapts []adapter.Adapter

	adapts = coingecko.GetTopHundred()

	if len(adapts) == 0 {
		slog.Error("GetExchange:coingecko.GetTopHundred", "message", "list is empty")
		return nil
	}

	var exchanges []model.ExchangeModel
	for i := 0; i < len(adapts); i++ {
		exchanges = append(exchanges, adapts[i].AdapterToExchange())
	}
	return exchanges
}


func GetDaily() []model.DailyModel {
	var adapts []adapter.Adapter
	adapts = coingecko.GetTopHundred()

	if len(adapts) == 0 {
		slog.Error("GetDaily:coingecko.GetTopHundred", "message", "list is empty")
		return nil
	} else {
		slog.Info("GetDaily:coingecko.GetTopHundred", "list", adapts)
	}

	var dailies []model.DailyModel
	for i := 0; i < len(adapts); i++ {
		dailies = append(dailies, adapts[i].AdapterToDaily(true))
	}
	return dailies
}
