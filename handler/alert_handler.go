package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/application/alert"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AlertHandler struct {
	Service *alert.Service
}

func NewAlertHandler(service *alert.Service) *AlertHandler {
	return &AlertHandler{
		Service: service,
	}
}

func (h *AlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Coin    string  `json:"coin"`
		Price   float32 `json:"price"`
		IsAbove bool    `json:"isAbove"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a, err := h.Service.CreateAlert(r.Context(), req.Coin, req.Price, req.IsAbove)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *AlertHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		Price   float32 `json:"price"`
		IsAbove bool    `json:"isAbove"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.Service.UpdateAlert(r.Context(), uint(id), req.Price, req.IsAbove); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AlertHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.Service.DeleteAlert(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AlertHandler) List(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.Service.ListAlerts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(alerts)
}
