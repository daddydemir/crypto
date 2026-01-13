package handler

import (
	"github.com/daddydemir/crypto/config/database"
	adiApp "github.com/daddydemir/crypto/pkg/analyses/adi/app"
	adiInfra "github.com/daddydemir/crypto/pkg/analyses/adi/infra"
	adiHandler "github.com/daddydemir/crypto/pkg/analyses/adi/rest"
	atrApp "github.com/daddydemir/crypto/pkg/analyses/atr/app"
	atrInfra "github.com/daddydemir/crypto/pkg/analyses/atr/infra"
	atrHandler "github.com/daddydemir/crypto/pkg/analyses/atr/rest"
	bollingerApp "github.com/daddydemir/crypto/pkg/analyses/bollinger/app"
	bollingerInfra "github.com/daddydemir/crypto/pkg/analyses/bollinger/infra"
	bollingerHandler "github.com/daddydemir/crypto/pkg/analyses/bollinger/rest"
	maApp "github.com/daddydemir/crypto/pkg/analyses/ma/app"
	maInfra "github.com/daddydemir/crypto/pkg/analyses/ma/infra"
	maHandler "github.com/daddydemir/crypto/pkg/analyses/ma/rest"
	rsiApp "github.com/daddydemir/crypto/pkg/analyses/rsi/app"
	rsiInfra "github.com/daddydemir/crypto/pkg/analyses/rsi/infra"
	rsiHandler "github.com/daddydemir/crypto/pkg/analyses/rsi/rest"

	emaApp "github.com/daddydemir/crypto/pkg/analyses/ema/app"
	emaInfra "github.com/daddydemir/crypto/pkg/analyses/ema/infra"
	emaHandler "github.com/daddydemir/crypto/pkg/analyses/ema/rest"
	"github.com/daddydemir/crypto/pkg/application/alert"

	binanceCandleApp "github.com/daddydemir/crypto/pkg/binance/application"
	binanceCandleInfra "github.com/daddydemir/crypto/pkg/binance/infrastructure"
	binanceCandleRest "github.com/daddydemir/crypto/pkg/binance/rest"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	donchianApp "github.com/daddydemir/crypto/pkg/channels/donchian/app"
	donchianInfra "github.com/daddydemir/crypto/pkg/channels/donchian/infra"
	donchianHandler "github.com/daddydemir/crypto/pkg/channels/donchian/rest"
	"github.com/daddydemir/crypto/pkg/factory"
	"github.com/daddydemir/crypto/pkg/infrastructure"

	"github.com/daddydemir/crypto/pkg/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

var serviceFactory *factory.ServiceFactory
var db = database.GetDatabaseService()
var cacheService = cache.GetCacheService()
var cacheable *service.CacheService

func init() {
	serviceFactory = factory.NewServiceFactory(db, cacheService, broker.GetBrokerService())
	cacheable = serviceFactory.NewCacheService()
}

func Route() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(setJSONContentType)
	r.Use(setLogging)

	r.HandleFunc("/health", healthHandler)

	base := "/api/v1"

	subRouter := r.PathPrefix(base).Subrouter()

	priceRepo := infrastructure.NewPriceRepository(cacheable, cacheService)

	//subRouter.HandleFunc("/topCoins", coinHandler.GetTopCoins).Methods(http.MethodGet)

	rsiHandler.NewHandler(rsiApp.NewApp(rsiInfra.NewRepository(cacheService, cacheable, db))).RegisterRoutes(subRouter)

	maHandler.NewHandler(maApp.NewApp(maInfra.NewRepository(cacheService), priceRepo)).RegisterRoutes(subRouter)

	emaHandler.NewHandler(emaApp.NewApp(emaInfra.NewRepository(cacheService))).RegisterRoutes(subRouter)

	bollingerHandler.NewHandler(bollingerApp.NewApp(bollingerInfra.NewRepository(cacheService), priceRepo)).RegisterRoutes(subRouter)

	alertHandler := NewAlertHandler(alert.NewService(infrastructure.NewAlertRepository(db)))

	subRouter.HandleFunc("/alerts", alertHandler.Create).Methods(http.MethodPost)
	subRouter.HandleFunc("/alerts/{id}", alertHandler.Update).Methods(http.MethodPut)
	subRouter.HandleFunc("/alerts/{id}", alertHandler.Delete).Methods(http.MethodDelete)
	subRouter.HandleFunc("/alerts", alertHandler.List).Methods(http.MethodGet)

	binanceCandleHandler := binanceCandleRest.NewCandleHandler(binanceCandleApp.NewGetCandlesQuery(binanceCandleInfra.NewCandleRepository(db)))
	subRouter.HandleFunc("/binance/coin/{symbol}", binanceCandleHandler.GetCandles).Methods(http.MethodGet)

	atrHandler.NewHandler(atrApp.NewApp(atrInfra.NewRepository(db))).RegisterRoutes(subRouter)

	donchianHandler.NewHandler(donchianApp.NewApp(donchianInfra.NewRepository(db))).RegisterRoutes(subRouter)

	adiHandler.NewHandler(adiApp.NewApp(adiInfra.NewRepository(db))).RegisterRoutes(subRouter)

	handler := cors.AllowAll().Handler(r)
	return handler
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
