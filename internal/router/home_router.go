package router

import (
	home "github.com/daddydemir/crypto/internal/handler/main"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/port/chart"
	serviceChart "github.com/daddydemir/crypto/internal/service/chart"
	"github.com/gorilla/mux"
)

func RegisterHomeRoutes(r *mux.Router, api port.CoinCapAPI, cache port.Cache) {

	rsiFactory := func(coin string) chart.RsiChart {
		return serviceChart.NewRsiService(coin, api, cache)
	}

	handler := home.NewHomeHandler(api, rsiFactory)

	r.HandleFunc("/", handler.Home)
}
