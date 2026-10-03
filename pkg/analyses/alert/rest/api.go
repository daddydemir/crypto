package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/alert/app"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	app *app.App
}

func username(r *http.Request) (string, bool) {
	value, ok := r.Context().Value("username").(string)
	return value, ok && value != ""
}

func NewHandler(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	u, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		Coin    string  `json:"coin"`
		Price   float32 `json:"price"`
		IsAbove bool    `json:"isAbove"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Coin == "" || req.Price <= 0 {
		http.Error(w, "coin and positive price are required", http.StatusBadRequest)
		return
	}
	a, err := h.app.CreateAlert(r.Context(), u, req.Coin, req.Price, req.IsAbove)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	u, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req struct {
		Price   float32 `json:"price"`
		IsAbove bool    `json:"isAbove"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Price <= 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	a, err := h.app.UpdateAlert(r.Context(), u, uint(id), req.Price, req.IsAbove)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	u, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.app.DeleteAlert(r.Context(), u, uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	u, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	alerts, err := h.app.ListAlerts(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(alerts)
}

func (h *Handler) SetStatus(w http.ResponseWriter, r *http.Request) {
	u, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req struct {
		IsActive bool `json:"isActive"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	a, err := h.app.SetStatus(r.Context(), u, uint(id), req.IsActive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/alerts", h.Create).Methods(http.MethodPost)
	router.HandleFunc("/alerts/{id}", h.Update).Methods(http.MethodPut)
	router.HandleFunc("/alerts/{id}", h.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/alerts/{id}/status", h.SetStatus).Methods(http.MethodPut)
	router.HandleFunc("/alerts", h.List).Methods(http.MethodGet)
}
