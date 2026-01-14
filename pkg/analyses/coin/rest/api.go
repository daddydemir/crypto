package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/coin/app"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) *Handler {
	return &Handler{
		app: app,
	}
}

func (h *Handler) TopCoins(w http.ResponseWriter, _ *http.Request) {
	coins, err := h.app.GetTopCoins()
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

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/topCoins", h.TopCoins)
}
