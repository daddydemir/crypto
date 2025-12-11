package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/application/bollinger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BollingerHandler struct {
	service *bollinger.Service
}

func NewBollingerHandler(service *bollinger.Service) *BollingerHandler {
	return &BollingerHandler{
		service: service,
	}
}

func (h *BollingerHandler) GetBollingerSeries(w http.ResponseWriter, r *http.Request) {
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

	series, err := h.service.GetBollingerSeries(r.Context(), coinID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(series)
}

func (h *BollingerHandler) BollingerBandSignals(w http.ResponseWriter, r *http.Request) {
	results, err := h.service.GetBollingerBandSignal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}
