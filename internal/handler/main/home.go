package home

import (
	"fmt"
	"github.com/daddydemir/crypto/assets"
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/port/chart"
	"log/slog"
	"net/http"
)

type HomeHandler struct {
	client         port.CoinCapAPI
	serviceFactory func(coin string) chart.RsiChart
}

func NewHomeHandler(c port.CoinCapAPI, factory func(coin string) chart.RsiChart) *HomeHandler {
	return &HomeHandler{client: c, serviceFactory: factory}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := assets.GetTemplate("templates/coin.html")

	err, coins := h.client.ListCoins()
	if err != nil {
		slog.Error("mainHandler:ListCoins", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//rsi := graphs.RSI{}
	var viewModels []model.CoinViewModel

	for i, coin := range coins {
		var class string
		var index float32 = 0

		if i <= 25 {
			factory := h.serviceFactory(coin.Id)
			index = factory.Index()
			//index = rsi.Index(coin.Id)
		}

		switch {
		case index == 0:
			class = ""
		case index >= 70:
			class = "bg-green-600 bg-opacity-50"
		case index <= 30:
			class = "bg-red-600 bg-opacity-50"
		default:
			class = "bg-yellow-400 bg-opacity-50"
		}

		viewModels = append(viewModels, model.CoinViewModel{
			Index:    i + 1,
			Name:     coin.Name,
			Symbol:   coin.Symbol,
			PriceUsd: coin.PriceUsd,
			Rsi:      index,
			RsiClass: class,
			Id:       coin.Id,
			GraphUrls: map[string]string{
				"rsi": fmt.Sprintf("/api/v1/graph/rsi/%s", coin.Id),
				"sma": fmt.Sprintf("/api/v1/graph/sma/%s", coin.Id),
				"ema": fmt.Sprintf("/api/v1/graph/ema/%s", coin.Id),
				"ma":  fmt.Sprintf("/api/v1/graph/ma/%s", coin.Id),
				"bb":  fmt.Sprintf("/api/v1/graph/bollingerBands/%s", coin.Id),
			},
		})
	}

	err = tmpl.Execute(w, struct {
		Coins []model.CoinViewModel
	}{
		Coins: viewModels,
	})
	if err != nil {
		slog.Error("mainHandler:tmpl.Execute", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
