package router

import (
	"github.com/daddydemir/crypto/internal/handler/alert"
	"github.com/daddydemir/crypto/internal/port"
	alertRepo "github.com/daddydemir/crypto/internal/port/alert"
	alertSvc "github.com/daddydemir/crypto/internal/service/alert"
	"github.com/daddydemir/crypto/internal/service/cache"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterAlertRoutes(r *mux.Router,
	repository alertRepo.AlertRepository,
	broker port.Broker,
	cache cache.CacheService,
) {

	service := alertSvc.NewAlertService(repository, broker)

	handler := alert.NewAlertHandler(*service, &cache)

	r.HandleFunc("/api/v1/alert", handler.ShowPage).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/alert", handler.Save).Methods(http.MethodPost)
}
