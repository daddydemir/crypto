package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/channels/donchian/app"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/donchian/coin/{symbol}", h.DonchianChannel)
}

func (h *Handler) DonchianChannel(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}

	series, err := h.app.Series(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(series) == 0 {
		http.Error(w, "no data found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(series)
}
