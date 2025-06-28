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

	RegisterGraphRoutes(r)
	RegisterDailyRoutes(r)
	RegisterExchangeRoutes(r)
	RegisterAlertRoutes(r)
	
	return cors.AllowAll().Handler(r)
}
