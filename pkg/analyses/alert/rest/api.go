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

func NewHandler(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Coin    string  `json:"coin"`
		Price   float32 `json:"price"`
		IsAbove bool    `json:"isAbove"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a, err := h.app.CreateAlert(r.Context(), req.Coin, req.Price, req.IsAbove)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		Price   float32 `json:"price"`
		IsAbove bool    `json:"isAbove"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.app.UpdateAlert(r.Context(), uint(id), req.Price, req.IsAbove); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.app.DeleteAlert(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.app.ListAlerts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(alerts)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/alerts", h.Create).Methods(http.MethodPost)
	router.HandleFunc("/alerts/{id}", h.Update).Methods(http.MethodPut)
	router.HandleFunc("/alerts/{id}", h.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/alerts", h.List).Methods(http.MethodGet)
}
