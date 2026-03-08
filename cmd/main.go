package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/daddydemir/crypto/config"
	cch "github.com/daddydemir/crypto/config/cache"
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/handler"
	"github.com/daddydemir/crypto/pkg/application/binance"
	"github.com/daddydemir/crypto/pkg/cache"
	binance2 "github.com/daddydemir/crypto/pkg/infrastructure/binance"
	"github.com/daddydemir/crypto/pkg/infrastructure/scheduler"
	bnnc "github.com/daddydemir/crypto/pkg/remote/binance"
	"github.com/daddydemir/crypto/pkg/service"
	_ "github.com/daddydemir/dlog"
	"gorm.io/gorm"
)

func main() {

	go bnnc.NewClient(config.Get("WS_URL"), cch.GetRedisClient()).Fetch()
	go handler.ListenAndBroadcast(cch.GetRedisClient())
	initJobs(database.GetDatabaseService())

	server := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              config.Get("PORT"),
		Handler:           handler.Route(),
	}

	if config.Get("ENV") == "PROD" {
		if err := server.ListenAndServeTLS(config.Get("CERT_PATH"), config.Get("KEY_PATH")); err != nil {
			slog.Error("ListenAndServeTLS", "error", err)
			panic(err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("ListenAndServe", "error", err)
			panic(err)
		}
	}

}

func initJobs(database *gorm.DB) {
	candleService := binance.NewCandleService(binance2.NewCandleRepository(database), binance2.NewDataSource())
	cacheService := service.NewCacheService(cache.GetCacheService())
	job := scheduler.FetchCandlesJob(candleService, *cacheService)

	job.Start()
}
