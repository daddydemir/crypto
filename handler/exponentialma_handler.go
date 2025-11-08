package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/application/exponentialma"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ExponentialMAHandler struct {
	service *exponentialma.Service
}

func NewExponentialMAHandler(service *exponentialma.Service) *ExponentialMAHandler {
	return &ExponentialMAHandler{
		service: service,
	}
}

func (h *ExponentialMAHandler) GetMovingAverages(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.GetMovingAverageSeries(coinID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
