package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/application/coin"
	"net/http"
)

type CoinHandler struct {
	usecase *coin.GetTopCoinsStats
	rsi     *coin.GetTopCoinsRSI
}

func NewCoinHandler(usecase *coin.GetTopCoinsStats, rsi *coin.GetTopCoinsRSI) *CoinHandler {
	return &CoinHandler{usecase: usecase, rsi: rsi}
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
