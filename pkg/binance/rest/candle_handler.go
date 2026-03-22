package rest

import (
	"encoding/json"
	"net/http"

	"github.com/daddydemir/crypto/pkg/binance/application"
	"github.com/gorilla/mux"
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
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	candles, err := h.getCandles.Execute(symbol, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(candles) == 0 {
		http.Error(w, "candle not found", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(candles)
	}
}
