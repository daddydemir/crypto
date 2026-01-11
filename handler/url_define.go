package handler

import (
	"github.com/daddydemir/crypto/config/database"
	adiApp "github.com/daddydemir/crypto/pkg/analyses/adi/app"
	adiInfra "github.com/daddydemir/crypto/pkg/analyses/adi/infra"
	adiRest "github.com/daddydemir/crypto/pkg/analyses/adi/rest"
	"github.com/daddydemir/crypto/pkg/application/alert"
	"github.com/daddydemir/crypto/pkg/application/bollinger"
	"github.com/daddydemir/crypto/pkg/application/coin"
	"github.com/daddydemir/crypto/pkg/application/exponentialma"
	"github.com/daddydemir/crypto/pkg/application/movingaverage"
	"github.com/daddydemir/crypto/pkg/atr/application"
	atrInfra "github.com/daddydemir/crypto/pkg/atr/infrastructure"
	"github.com/daddydemir/crypto/pkg/atr/rest"
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
	coinInfra "github.com/daddydemir/crypto/pkg/infrastructure/coin"
	expoInfra "github.com/daddydemir/crypto/pkg/infrastructure/exponentialma"
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

	usecase := coin.NewGetTopCoinsStats(coinInfra.NewCacheHistoryRepository(cacheService), coinInfra.NewCoinGeckoMarketRepository(serviceFactory.NewCachedCoinCapClient(), db))
	rsi := coin.NewGetTopCoinsRSI(coinInfra.NewPriceRepository(cacheService, cacheable, db))
	rsiHistory := coin.NewGetCoinRSIHistory(coinInfra.NewPriceRepository(cacheService, cacheable, db))
	coinHandler := NewCoinHandler(usecase, rsi, rsiHistory)

	infraPriceRepository := infrastructure.NewPriceRepository(cacheable, cacheService)
	movingAverageHandler := NewMovingAverageHandler(movingaverage.NewService(infrastructure.NewPriceHistoryRepository(cacheService), infraPriceRepository))
	exponentialHandler := NewExponentialMAHandler(exponentialma.NewService(expoInfra.NewPriceHistoryRepository(cacheService)))

	bollingerHandler := NewBollingerHandler(bollinger.NewService(infrastructure.NewBollingerRepository(cacheService), infraPriceRepository))

	subRouter.HandleFunc("/topCoins", coinHandler.GetTopCoins).Methods(http.MethodGet)
	subRouter.HandleFunc("/topCoinsRSI", coinHandler.GetTopCoinsRSI).Methods(http.MethodGet)
	subRouter.HandleFunc("/coins/{id}/rsi/history", coinHandler.GetCoinRSIHistory).Methods(http.MethodGet)
	subRouter.HandleFunc("/coins/moving-averages", movingAverageHandler.MovingAverageSignals).Methods(http.MethodGet)
	subRouter.HandleFunc("/coins/{id}/moving-averages", movingAverageHandler.GetMovingAverages).Methods(http.MethodGet)
	subRouter.HandleFunc("/coins/{id}/exponential-moving-averages", exponentialHandler.GetMovingAverages).Methods(http.MethodGet)
	subRouter.HandleFunc("/coins/{id}/bollinger-bands", bollingerHandler.GetBollingerSeries).Methods(http.MethodGet)
	subRouter.HandleFunc("/coins/bollinger-bands", bollingerHandler.BollingerBandSignals).Methods(http.MethodGet)

	alertHandler := NewAlertHandler(alert.NewService(infrastructure.NewAlertRepository(db)))

	subRouter.HandleFunc("/alerts", alertHandler.Create).Methods(http.MethodPost)
	subRouter.HandleFunc("/alerts/{id}", alertHandler.Update).Methods(http.MethodPut)
	subRouter.HandleFunc("/alerts/{id}", alertHandler.Delete).Methods(http.MethodDelete)
	subRouter.HandleFunc("/alerts", alertHandler.List).Methods(http.MethodGet)

	binanceCandleHandler := binanceCandleRest.NewCandleHandler(binanceCandleApp.NewGetCandlesQuery(binanceCandleInfra.NewCandleRepository(db)))
	subRouter.HandleFunc("/binance/coin/{symbol}", binanceCandleHandler.GetCandles).Methods(http.MethodGet)

	atrHandler := rest.NewAtrHandler(application.NewPointService(atrInfra.NewAtrRepository(db)))
	subRouter.HandleFunc("/atr/coin/{symbol}", atrHandler.Points).Methods(http.MethodGet)

	donchianHandler.NewDonchianHandler(donchianApp.NewDonchianApp(donchianInfra.NewDonchianRepository(db))).RegisterRoutes(subRouter)
	adiRest.NewHandler(adiApp.NewApp(adiInfra.NewRepository(db))).RegisterRoutes(subRouter)

	handler := cors.AllowAll().Handler(r)
	return handler
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
