package handler

import (
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/factory"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

var serviceFactory *factory.ServiceFactory

func init() {
	serviceFactory = factory.NewServiceFactory(database.GetDatabaseService(), cache.GetCacheService(), broker.GetBrokerService())
	alertService = serviceFactory.NewAlertService()
	cacheService = serviceFactory.NewCacheService()
}

func Route() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(setJSONContentType)
	r.Use(setLogging)

	base := "/api/v1"

	subRouter := r.PathPrefix(base).Subrouter()
	subRouter.HandleFunc("/graph/rsi/{coin}", rsiHandler).Methods(http.MethodGet)
	subRouter.HandleFunc("/graph/sma/{coin}", smaHandler).Methods(http.MethodGet)
	subRouter.HandleFunc("/graph/ema/{coin}", emaHandler).Methods(http.MethodGet)
	subRouter.HandleFunc("/graph/ma/{coin}", maHandler).Methods(http.MethodGet)
	subRouter.HandleFunc("/graph/bollingerBands/{coin}", bollingerBandsHandler).Methods(http.MethodGet)
	subRouter.HandleFunc("/graph/main", mainHandler).Methods(http.MethodGet)

	subRouter.HandleFunc("/dailyStart", dailyStart).Methods(http.MethodGet)
	subRouter.HandleFunc("/dailyEnd", dailyEnd).Methods(http.MethodGet)
	subRouter.HandleFunc("/daily", daily).Methods(http.MethodGet)
	subRouter.HandleFunc("/getDaily", getDaily).Methods(http.MethodPost)
	subRouter.HandleFunc("/getDailyWithId", getDailyWithId).Methods(http.MethodPost)

	subRouter.HandleFunc("/exchange", getExchange).Methods(http.MethodGet)
	subRouter.HandleFunc("/getExchange", getExchangeFromDb).Methods(http.MethodGet)
	subRouter.HandleFunc("/createExchange", createExchange).Methods(http.MethodGet)

	subRouter.HandleFunc("/weekly", getWeekly).Methods(http.MethodGet)

	subRouter.HandleFunc("/alert", alertPage).Methods(http.MethodGet)
	subRouter.HandleFunc("/alert", alert).Methods(http.MethodPost)

	handler := cors.AllowAll().Handler(r)
	return handler
}
