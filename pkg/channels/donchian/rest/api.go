package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/daddydemir/crypto/pkg/channels/donchian/app"
	"github.com/gorilla/mux"
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

	periodStr := r.URL.Query().Get("period")

	period := 20
	if periodStr != "" {
		period, _ = strconv.Atoi(periodStr)
	}

	series, err := h.app.Series(symbol, period)
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
