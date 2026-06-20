package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/daddydemir/crypto/pkg/analyses/ma/app"
	"github.com/gorilla/mux"
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
	shortStr := r.URL.Query().Get("short")
	midStr := r.URL.Query().Get("mid")
	longStr := r.URL.Query().Get("long")

	days := 99
	short, mid, long := 7, 25, 99
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}
	if shortStr != "" {
		short, _ = strconv.Atoi(shortStr)
	}
	if midStr != "" {
		mid, _ = strconv.Atoi(midStr)
	}
	if longStr != "" {
		long, _ = strconv.Atoi(longStr)
	}

	result, err := h.app.GetMovingAverageSeries(coinID, days, short, mid, long)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) MovingAverageSignals(w http.ResponseWriter, _ *http.Request) {
	results, err := h.app.GetMovingAverageSignals()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/coins/{id}/moving-averages", h.MovingAverages)
	router.HandleFunc("/coins/moving-averages", h.MovingAverageSignals)
}
