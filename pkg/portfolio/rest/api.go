package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/daddydemir/crypto/pkg/portfolio/app"
	"github.com/daddydemir/crypto/pkg/portfolio/domain"
	"github.com/daddydemir/crypto/pkg/portfolio/exchange"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	app      *app.App
	exchange *exchange.Client
}

func NewHandler(application *app.App, exchangeClient *exchange.Client) *Handler {
	return &Handler{app: application, exchange: exchangeClient}
}

func (h *Handler) FetchExchangeTrades(w http.ResponseWriter, r *http.Request) {
	if _, ok := username(r); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var request struct {
		Exchange, BaseAsset, QuoteAsset string
		StartDate, EndDate              time.Time
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Geçersiz sorgu bilgileri")
		return
	}
	transactions, err := h.exchange.Fetch(r.Context(), request.Exchange, request.BaseAsset, request.QuoteAsset, request.StartDate, request.EndDate)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func (h *Handler) ImportExchangeTrades(w http.ResponseWriter, r *http.Request) {
	user, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var transactions []domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transactions); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	count, err := h.app.Import(user, transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"imported": count})
}

func username(r *http.Request) (string, bool) {
	value, ok := r.Context().Value("username").(string)
	return value, ok && value != ""
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}

func transactionID(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	return uint(id), err
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	transactions, err := h.app.List(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var transaction domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.app.Create(user, transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := transactionID(r)
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}
	var transaction domain.Transaction
	if err = json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updated, err := h.app.Update(user, id, transaction)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "portfolio transaction not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := username(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := transactionID(r)
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}
	if err = h.app.Delete(user, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/portfolio/transactions", h.List).Methods(http.MethodGet)
	router.HandleFunc("/portfolio/transactions", h.Create).Methods(http.MethodPost)
	router.HandleFunc("/portfolio/transactions/{id}", h.Update).Methods(http.MethodPut)
	router.HandleFunc("/portfolio/transactions/{id}", h.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/portfolio/exchange-trades/search", h.FetchExchangeTrades).Methods(http.MethodPost)
	router.HandleFunc("/portfolio/exchange-trades/import", h.ImportExchangeTrades).Methods(http.MethodPost)
}
