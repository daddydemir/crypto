package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/application/movingaverage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type MovingAverageHandler struct {
	Service *movingaverage.Service
}

func NewMovingAverageHandler(service *movingaverage.Service) *MovingAverageHandler {
	return &MovingAverageHandler{
		Service: service,
	}
}

func (h *MovingAverageHandler) GetMovingAverages(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.Service.GetMovingAverageSeries(r.Context(), coinID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *MovingAverageHandler) MovingAverageSignals(w http.ResponseWriter, r *http.Request) {
	results, err := h.Service.GetMovingAverageSignals()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}
