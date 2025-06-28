package main

import (
	"github.com/daddydemir/crypto/config/broker"
	"github.com/daddydemir/crypto/config/database"
	brokeradapter "github.com/daddydemir/crypto/internal/adapter/broker"
	cacheadapter "github.com/daddydemir/crypto/internal/adapter/cache"
	coincapadapter "github.com/daddydemir/crypto/internal/adapter/coincap"
	"github.com/daddydemir/crypto/internal/adapter/coingecko"
	"github.com/daddydemir/crypto/internal/adapter/repository/postgres"
	"github.com/daddydemir/crypto/internal/config"
	cronjob "github.com/daddydemir/crypto/internal/cron"
	"github.com/daddydemir/crypto/internal/cron/jobs"
	"github.com/daddydemir/crypto/internal/router"
	"github.com/daddydemir/crypto/internal/service/cache"
	"github.com/daddydemir/crypto/internal/service/chart"
	"github.com/daddydemir/crypto/internal/service/daily"
	"github.com/daddydemir/crypto/internal/service/maintenance"
	"github.com/daddydemir/crypto/internal/service/movingaverage"
	"github.com/daddydemir/crypto/internal/service/signal"
	//_ "github.com/daddydemir/crypto/pkg/cronjob"
	"github.com/daddydemir/crypto/pkg/remote/coincap/client"
	"github.com/daddydemir/crypto/pkg/token/provider"
	"github.com/daddydemir/crypto/pkg/token/strategy"
	_ "github.com/daddydemir/dlog"
	"log"
	"net/http"
)

func main() {

	db := database.GetDatabaseService()
	cacheService := cacheadapter.NewRedisCache(config.NewRedisClient())

	dailyService := daily.NewDefaultDailyService(postgres.NewPostgresDailyRepository(db))
	dailyCreator := daily.NewDailyCreator(dailyService, coingecko.NewHttpGeckoClient(&http.Client{}))

	coinCapClient := coincapadapter.NewRealCoinCapClient("https://rest.coincap.io", clientFactory)

	publisher, _ := brokeradapter.NewRabbitPublisher(broker.GetChannel(), "assistant")

	ss := signal.NewDefaultSignalService(postgres.NewPostgresSignalRepository(db))

	rsiService := chart.NewRsiService("", coinCapClient, cacheService)
	movingAverageService := movingaverage.NewService(cacheService, ss, publisher, rsiService)

	// todo: add this services
	scheduler := cronjob.NewScheduler(
		jobs.NewDailyStartJob(dailyService, *dailyCreator),
		jobs.NewDailyEndJob(dailyService, *dailyCreator),
		jobs.NewValidateCacheJob(maintenance.NewValidateService(cacheService, coinCapClient)),
		jobs.NewCheckAllMAJob(*movingAverageService),
	)

	scheduler.Start()

	cachedCoinCapClient := coincapadapter.NewCachedCoinCapClient(coinCapClient, cacheService)

	repository := postgres.NewPostgresAlertRepository(db)

	cacheImpl := cache.NewCacheService(cacheService)

	router := router.NewRouter(
		cacheService,
		cachedCoinCapClient,
		repository,
		publisher,
		*cacheImpl,
		db,
	)
	log.Fatal(http.ListenAndServe(":8080", router))

	//server := &http.Server{
	//	Addr:    config.Get("PORT"),
	//	Handler: handler.Route(),
	//}
	//
	//if config.Get("ENV") == "PROD" {
	//	if err := server.ListenAndServeTLS(config.Get("CERT_PATH"), config.Get("KEY_PATH")); err != nil {
	//		slog.Error("ListenAndServeTLS", "error", err)
	//		panic(err)
	//	}
	//} else {
	//	if err := server.ListenAndServe(); err != nil {
	//		slog.Error("ListenAndServe", "error", err)
	//		panic(err)
	//	}
	//}

}

func clientFactory(ts coincapadapter.TokenStrategy) client.TokenAwareClient {
	switch ts {
	case coincapadapter.QueryStrategy:
		return *client.NewTokenAwareClient(
			"https://rest.coincap.io",
			provider.NewRedisTokenProvider(),
			strategy.QueryTokenStrategy{ParamName: "apiKey"},
		)
	case coincapadapter.HeaderStrategy:
		return *client.NewTokenAwareClient(
			"https://rest.coincap.io",
			provider.NewRedisTokenProvider(),
			strategy.HeaderTokenStrategy{},
		)
	default:
		panic("unsupported token strategy")
	}
}

func startApp() {

	realClient := coincapadapter.NewRealCoinCapClient("https://rest.coincap.io", clientFactory)
	realClient.ListCoins()

}
