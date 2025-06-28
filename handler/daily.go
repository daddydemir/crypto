package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterDailyRoutes(router *mux.Router) {

	router.HandleFunc("/api/v1/dailyStart", dailyStart).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/dailyEnd", dailyEnd).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/daily", daily).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/getDaily", getDaily).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/getDailyWithId", getDailyWithId).Methods(http.MethodPost)

}
