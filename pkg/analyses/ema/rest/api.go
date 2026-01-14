package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/ema/app"
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

func (h *Handler) MovingAverages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coinID := vars["id"]

	if coinID == "" {
		http.Error(w, "coin param required", http.StatusBadRequest)
		return
	}

	daysStr := r.URL.Query().Get("days")

	days := 99
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	result, err := h.app.GetMovingAverageSeries(coinID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/coins/{id}/exponential-moving-averages", h.MovingAverages)
}
