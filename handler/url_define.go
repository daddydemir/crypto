package handler

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func Route() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(setJSONContentType)

	base := "/api/v1"

	subRouter := r.PathPrefix(base).Subrouter()
	subRouter.HandleFunc("/graph/rsi/{coin}", rsiHandler).Methods(http.MethodGet)

	subRouter.HandleFunc("/dailyStart", dailyStart).Methods(http.MethodGet)
	subRouter.HandleFunc("/dailyEnd", dailyEnd).Methods(http.MethodGet)
	subRouter.HandleFunc("/daily", daily).Methods(http.MethodGet)
	subRouter.HandleFunc("/getDaily", getDaily).Methods(http.MethodPost)
	subRouter.HandleFunc("/getDailyWithId", getDailyWithId).Methods(http.MethodPost)

	subRouter.HandleFunc("/exchange", getExchange).Methods(http.MethodGet)
	subRouter.HandleFunc("/getExchange", getExchangeFromDb).Methods(http.MethodGet)
	subRouter.HandleFunc("/createExchange", createExchange).Methods(http.MethodGet)

	subRouter.HandleFunc("/weekly", getWeekly).Methods(http.MethodGet)

	handler := cors.AllowAll().Handler(r)
	return handler
}
