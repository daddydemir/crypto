package router

import (
	"github.com/daddydemir/crypto/internal/handler/daily"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterDailyRoutes(r *mux.Router, handler *daily.DailyHandler) {

	r.HandleFunc("/api/v1/daily", handler.GetByDateRange).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/daily/{id}", handler.GetByIdAndDateRange).Methods(http.MethodPost)

}
