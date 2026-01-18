package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/macd/app"
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

func (h *Handler) CoinMACDHistory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	coinSymbol := vars["symbol"]

	data, err := h.app.CoinMACDHistory(coinSymbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/macd/coin/{symbol}", h.CoinMACDHistory).Methods(http.MethodGet)
}
