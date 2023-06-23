package handler

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func Route() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(setJSONContentType)

	base := "/api/v1/"

	r.HandleFunc(base+"dailyStart", dailyStart).Methods(http.MethodGet)
	r.HandleFunc(base+"dailyEnd", dailyEnd).Methods(http.MethodGet)
	r.HandleFunc(base+"daily", daily).Methods(http.MethodGet)
	r.HandleFunc(base+"getDaily", getDaily).Methods(http.MethodPost)

	r.HandleFunc(base+"exchange", getExchange).Methods(http.MethodGet)
	r.HandleFunc(base+"getExchange", getExchangeFromDb).Methods(http.MethodGet)
	r.HandleFunc(base+"createExchange", createExchange).Methods(http.MethodGet)

	r.HandleFunc(base+"weekly", getWeekly).Methods(http.MethodGet)

	handler := cors.AllowAll().Handler(r)
	return handler
}
