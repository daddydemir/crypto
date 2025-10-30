package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/application/coin"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type CoinHandler struct {
	usecase    *coin.GetTopCoinsStats
	rsi        *coin.GetTopCoinsRSI
	rsiHistory *coin.GetCoinRSIHistory
}

func NewCoinHandler(usecase *coin.GetTopCoinsStats, rsi *coin.GetTopCoinsRSI, rsiHistory *coin.GetCoinRSIHistory) *CoinHandler {
	return &CoinHandler{usecase: usecase, rsi: rsi, rsiHistory: rsiHistory}
}

func (h *CoinHandler) GetTopCoins(w http.ResponseWriter, _ *http.Request) {
	coins, err := h.usecase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(coins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CoinHandler) GetTopCoinsRSI(w http.ResponseWriter, _ *http.Request) {
	data, err := h.rsi.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CoinHandler) GetCoinRSIHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coinID := vars["id"]
	daysStr := r.URL.Query().Get("days")

	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	data, err := h.rsiHistory.Execute(coinID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}
