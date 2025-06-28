package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterAlertRoutes(router *mux.Router) {

	router.HandleFunc("/api/v1/alert", alertPage).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/alert", alert).Methods(http.MethodPost)
}
