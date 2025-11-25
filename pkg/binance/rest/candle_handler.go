package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/binance/application"
	"github.com/gorilla/mux"
	"net/http"
)

type CandleHandler struct {
	getCandles *application.GetCandlesQuery
}

func NewCandleHandler(application *application.GetCandlesQuery) *CandleHandler {
	return &CandleHandler{getCandles: application}
}

func (h *CandleHandler) GetCandles(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
	}
	candles, err := h.getCandles.Execute(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(candles) == 0 {
		http.Error(w, "candle not found", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(candles)
	}
}
