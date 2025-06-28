package router

import (
	handlerChart "github.com/daddydemir/crypto/internal/handler/chart"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/port/chart"
	serviceChart "github.com/daddydemir/crypto/internal/service/chart"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterGraphRoutes(r *mux.Router, cache port.Cache, client port.CoinCapAPI) {
	base := "/api/v1/graph"

	// Factory: coin'e göre handler'a BollingerService döner
	bollingerFactory := func(coin string) chart.BollingerChart {
		return serviceChart.NewBollingerBandsService(coin, 20, cache)
	}

	emaFactory := func(coin string) chart.EmaChart {
		return serviceChart.NewEmaService(coin, 25, cache)
	}

	smaFactory := func(coin string) chart.SmaChart {
		return serviceChart.NewSmaService(coin, 10, cache)
	}

	tripleEmaFactory := func(coin string) chart.TripleEmaChart {
		return serviceChart.NewTripleEmaService(coin, cache)
	}

	rsiFactory := func(coin string) chart.RsiChart {
		return serviceChart.NewRsiService(coin, client, cache)
	}

	bollingerHandler := handlerChart.NewBollingerHandler(bollingerFactory)
	r.HandleFunc(base+"/bollingerBands/{coin}", bollingerHandler.Handle).Methods(http.MethodGet)

	emaHandler := handlerChart.NewEmaHandler(emaFactory)
	r.HandleFunc(base+"/ema/{coin}", emaHandler.Handle).Methods(http.MethodGet)

	smaHandler := handlerChart.NewSmaHandler(smaFactory)
	r.HandleFunc(base+"/sma/{coin}", smaHandler.Handle).Methods(http.MethodGet)

	tripleEmaHandler := handlerChart.NewTripleEmaHandler(tripleEmaFactory)
	r.HandleFunc(base+"/ma/{coin}", tripleEmaHandler.Handle).Methods(http.MethodGet)

	rsiHandler := handlerChart.NewRsiHandler(rsiFactory)
	r.HandleFunc("/api/v1/graph/rsi/{coin}", rsiHandler.Handle).Methods(http.MethodGet)
}
