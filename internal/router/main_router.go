package router

import (
	"github.com/daddydemir/crypto/internal/adapter/repository/postgres"
	dh "github.com/daddydemir/crypto/internal/handler/daily"
	"github.com/daddydemir/crypto/internal/middleware"
	"github.com/daddydemir/crypto/internal/port"
	alertRepo "github.com/daddydemir/crypto/internal/port/alert"
	"github.com/daddydemir/crypto/internal/service/cache"
	"github.com/daddydemir/crypto/internal/service/daily"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"net/http"
)

func NewRouter(
	cache port.Cache,
	api port.CoinCapAPI,
	repository alertRepo.AlertRepository,
	broker port.Broker,
	cacheSvc cache.CacheService,
	db *gorm.DB,
) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(middleware.Logging)
	r.Use(middleware.JSONContentType)

	// Modüler route kayıtları
	RegisterGraphRoutes(r, cache, api)
	RegisterAlertRoutes(r, repository, broker, cacheSvc)

	dailyHandler := dh.NewDailyHandler(daily.NewDefaultDailyService(postgres.NewPostgresDailyRepository(db)))
	RegisterDailyRoutes(r, dailyHandler)
	RegisterHomeRoutes(r, api, cache)

	// RegisterAlertRoutes(r, ...)
	// RegisterExchangeRoutes(r, ...)

	return cors.AllowAll().Handler(r)
}
