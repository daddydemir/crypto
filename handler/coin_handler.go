package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/application/coin"
	"net/http"
)

type CoinHandler struct {
	usecase *coin.GetTopCoinsStats
}

func NewCoinHandler(usecase *coin.GetTopCoinsStats) *CoinHandler {
	return &CoinHandler{usecase: usecase}
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
