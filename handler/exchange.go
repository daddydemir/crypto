package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterExchangeRoutes(router *mux.Router) {

	router.HandleFunc("/api/v1/exchange", getExchange).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/getExchange", getExchangeFromDb).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/createExchange", createExchange).Methods(http.MethodGet)
}
