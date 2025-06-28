package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterGraphRoutes(router *mux.Router) {

	router.HandleFunc("/api/v1/graph/rsi/{coin}", rsiHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/graph/sma/{coin}", smaHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/graph/ema/{coin}", emaHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/graph/ma/{coin}", maHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/graph/bollingerBands/{coin}", bollingerBandsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/graph/main", mainHandler).Methods(http.MethodGet)

	router.Handle("/api/v2/graph/bollingerBands/{coin}", nil)

}
