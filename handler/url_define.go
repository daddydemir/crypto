package handler

import (
	"log/slog"

	"github.com/daddydemir/crypto/config"
	redisConfig "github.com/daddydemir/crypto/config/cache"
	"github.com/daddydemir/crypto/config/database"
	adiApp "github.com/daddydemir/crypto/pkg/analyses/adi/app"
	adiInfra "github.com/daddydemir/crypto/pkg/analyses/adi/infra"
	adiHandler "github.com/daddydemir/crypto/pkg/analyses/adi/rest"
	alertApp "github.com/daddydemir/crypto/pkg/analyses/alert/app"
	alertInfra "github.com/daddydemir/crypto/pkg/analyses/alert/infra"
	alertHandler "github.com/daddydemir/crypto/pkg/analyses/alert/rest"
	atrApp "github.com/daddydemir/crypto/pkg/analyses/atr/app"
	atrInfra "github.com/daddydemir/crypto/pkg/analyses/atr/infra"
	atrHandler "github.com/daddydemir/crypto/pkg/analyses/atr/rest"
	bollingerApp "github.com/daddydemir/crypto/pkg/analyses/bollinger/app"
	bollingerInfra "github.com/daddydemir/crypto/pkg/analyses/bollinger/infra"
	bollingerHandler "github.com/daddydemir/crypto/pkg/analyses/bollinger/rest"
	coinApp "github.com/daddydemir/crypto/pkg/analyses/coin/app"
	coinInfra "github.com/daddydemir/crypto/pkg/analyses/coin/infra"
	coinHandler "github.com/daddydemir/crypto/pkg/analyses/coin/rest"
	maApp "github.com/daddydemir/crypto/pkg/analyses/ma/app"
	maInfra "github.com/daddydemir/crypto/pkg/analyses/ma/infra"
	maHandler "github.com/daddydemir/crypto/pkg/analyses/ma/rest"
	macdApp "github.com/daddydemir/crypto/pkg/analyses/macd/app"
	macdInfra "github.com/daddydemir/crypto/pkg/analyses/macd/infra"
	macdHandler "github.com/daddydemir/crypto/pkg/analyses/macd/rest"
	marketbreadthApp "github.com/daddydemir/crypto/pkg/analyses/marketbreadth/app"
	marketbreadthInfra "github.com/daddydemir/crypto/pkg/analyses/marketbreadth/infra"
	marketbreadthHandler "github.com/daddydemir/crypto/pkg/analyses/marketbreadth/rest"
	notfyApp "github.com/daddydemir/crypto/pkg/analyses/notification/app"
	notfyInfra "github.com/daddydemir/crypto/pkg/analyses/notification/infra"
	notfyHandler "github.com/daddydemir/crypto/pkg/analyses/notification/rest"
	rsiApp "github.com/daddydemir/crypto/pkg/analyses/rsi/app"
	rsiInfra "github.com/daddydemir/crypto/pkg/analyses/rsi/infra"
	rsiHandler "github.com/daddydemir/crypto/pkg/analyses/rsi/rest"
	basicApp "github.com/daddydemir/crypto/pkg/auth/basic/app"
	basicInfra "github.com/daddydemir/crypto/pkg/auth/basic/infra"
	basicHandler "github.com/daddydemir/crypto/pkg/auth/basic/rest"
	"github.com/daddydemir/crypto/pkg/token/jwt"

	emaApp "github.com/daddydemir/crypto/pkg/analyses/ema/app"
	emaInfra "github.com/daddydemir/crypto/pkg/analyses/ema/infra"
	emaHandler "github.com/daddydemir/crypto/pkg/analyses/ema/rest"

	binanceCandleApp "github.com/daddydemir/crypto/pkg/binance/application"
	binanceCandleInfra "github.com/daddydemir/crypto/pkg/binance/infrastructure"
	binanceCandleRest "github.com/daddydemir/crypto/pkg/binance/rest"
	"github.com/daddydemir/crypto/pkg/cache"
	donchianApp "github.com/daddydemir/crypto/pkg/channels/donchian/app"
	donchianInfra "github.com/daddydemir/crypto/pkg/channels/donchian/infra"
	donchianHandler "github.com/daddydemir/crypto/pkg/channels/donchian/rest"
	globalSearch "github.com/daddydemir/crypto/pkg/globalsearch"
	"github.com/daddydemir/crypto/pkg/infrastructure"
	portfolioApp "github.com/daddydemir/crypto/pkg/portfolio/app"
	portfolioExchange "github.com/daddydemir/crypto/pkg/portfolio/exchange"
	portfolioInfra "github.com/daddydemir/crypto/pkg/portfolio/infra"
	portfolioHandler "github.com/daddydemir/crypto/pkg/portfolio/rest"
	strategyApp "github.com/daddydemir/crypto/pkg/strategylab/app"
	strategyInfra "github.com/daddydemir/crypto/pkg/strategylab/infra"
	strategyHandler "github.com/daddydemir/crypto/pkg/strategylab/rest"
	strategyScheduler "github.com/daddydemir/crypto/pkg/strategylab/scheduler"
	tradeExplorer "github.com/daddydemir/crypto/pkg/tradeexplorer"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db = database.GetDatabaseService()
var cacheService = cache.GetCacheService()
var tokenService = jwt.NewTokenService(config.Get("JWT_SECRET"))

func Route() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(setJSONContentType)
	r.Use(setLogging)

	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/ws", handleConnections)

	base := "/api/v1"

	subRouter := r.PathPrefix(base).Subrouter()

	authorize := r.PathPrefix(base).Subrouter()
	authorize.Use(auth)

	coinCatalog := infrastructure.NewCoinCatalog(db, cacheService)
	if _, err := coinCatalog.List(); err != nil {
		slog.Error("failed to initialize coin catalog cache", "error", err)
	}
	priceRepo := infrastructure.NewPriceRepository(db, coinCatalog)

	coinHandler.NewHandler(coinApp.NewApp(coinInfra.NewRepository(db, coinCatalog))).RegisterRoutes(subRouter)

	rsiHandler.NewHandler(rsiApp.NewApp(rsiInfra.NewRepository(db, coinCatalog))).RegisterRoutes(subRouter)

	maHandler.NewHandler(maApp.NewApp(maInfra.NewRepository(db), priceRepo)).RegisterRoutes(subRouter)
	marketbreadthHandler.NewHandler(marketbreadthApp.NewApp(marketbreadthInfra.NewRepository(db), coinCatalog, cacheService)).RegisterRoutes(subRouter)

	emaHandler.NewHandler(emaApp.NewApp(emaInfra.NewRepository(db))).RegisterRoutes(subRouter)

	bollingerHandler.NewHandler(bollingerApp.NewApp(bollingerInfra.NewRepository(db), priceRepo)).RegisterRoutes(subRouter)

	alertHandler.NewHandler(alertApp.NewApp(alertInfra.NewRepository(db))).RegisterRoutes(authorize)
	portfolioRepository := portfolioInfra.NewRepository(db)
	if err := portfolioRepository.Migrate(); err != nil {
		slog.Error("failed to migrate portfolio transactions", "error", err)
	}
	portfolioHandler.NewHandler(portfolioApp.NewApp(portfolioRepository), portfolioExchange.NewClient()).RegisterRoutes(authorize)
	strategyRepository := strategyInfra.NewRepository(db)
	if err := strategyRepository.Migrate(); err != nil {
		slog.Error("failed to migrate strategy lab", "error", err)
	}
	strategyApplication := strategyApp.New(strategyRepository, redisConfig.GetRedisClient())
	strategyHandler.New(strategyApplication).Register(authorize)
	strategyApplication.StartBacktestWorker()
	strategyScheduler.Start(strategyApplication)
	tradeExplorer.NewHandler(tradeExplorer.New(db, redisConfig.GetRedisClient())).Register(authorize)
	globalSearch.NewHandler(globalSearch.New(db, coinCatalog)).Register(authorize)

	binanceCandleHandler := binanceCandleRest.NewCandleHandler(binanceCandleApp.NewGetCandlesQuery(binanceCandleInfra.NewCandleRepository(db)))
	subRouter.HandleFunc("/binance/coin/{symbol}", binanceCandleHandler.GetCandles).Methods(http.MethodGet)

	atrHandler.NewHandler(atrApp.NewApp(atrInfra.NewRepository(db))).RegisterRoutes(subRouter)

	donchianHandler.NewHandler(donchianApp.NewApp(donchianInfra.NewRepository(db))).RegisterRoutes(subRouter)

	adiHandler.NewHandler(adiApp.NewApp(adiInfra.NewRepository(db))).RegisterRoutes(subRouter)

	macdHandler.NewHandler(macdApp.NewApp(macdInfra.NewRepository(db))).RegisterRoutes(subRouter)

	notfyHandler.NewHandler(notfyApp.NewApp(notfyInfra.NewRepository(cacheService))).RegisterRoutes(authorize)
	basicHandler.NewHandler(basicApp.NewApp(basicInfra.NewRepository(db, tokenService))).RegisterRoutes(subRouter)

	handler := cors.AllowAll().Handler(r)
	return handler
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
