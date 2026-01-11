package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/bollinger/app"
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

func (h *Handler) BollingerSeries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coinID := vars["id"]

	if coinID == "" {
		http.Error(w, "coin param required", http.StatusBadRequest)
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 0
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	series, err := h.app.GetBollingerSeries(coinID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(series)
}

func (h *Handler) BollingerBandSignals(w http.ResponseWriter, _ *http.Request) {
	results, err := h.app.GetBollingerBandSignal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/coins/{id}/bollinger-bands", h.BollingerSeries)
	router.HandleFunc("/bollinger-bands", h.BollingerBandSignals)
}
