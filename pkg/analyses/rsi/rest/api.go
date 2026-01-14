package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/rsi/app"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) *Handler {
	return &Handler{
		app: app,
	}
}

//func (h *Handler) TopCoins(w http.ResponseWriter, _ *http.Request) {
//	coins, err := h.usecase.Execute()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = json.NewEncoder(w).Encode(coins)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

func (h *Handler) TopCoinsRSI(w http.ResponseWriter, _ *http.Request) {
	data, err := h.app.TopCoinsRSI()
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

func (h *Handler) CoinHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coinID := vars["id"]
	daysStr := r.URL.Query().Get("days")

	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	data, err := h.app.CoinRSIHistory(coinID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/topCoinsRSI", h.TopCoinsRSI)
	router.HandleFunc("/coins/{id}/rsi/history", h.CoinHistory)
}
