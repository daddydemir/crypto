package alert

import (
	"encoding/json"
	"github.com/daddydemir/crypto/internal/port/cache"
	"net/http"
	"strconv"

	"github.com/daddydemir/crypto/assets"
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/service/alert"
	"github.com/gorilla/mux"
	"log/slog"
)

type AlertHandler struct {
	service alert.AlertService
	cache   cache.CoinCacheService
}

func NewAlertHandler(service alert.AlertService, cache cache.CoinCacheService) *AlertHandler {
	return &AlertHandler{
		service: service,
		cache:   cache,
	}
}

func (h *AlertHandler) ShowPage(w http.ResponseWriter, r *http.Request) {
	tmpl := assets.GetTemplate("templates/alert.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	coins := h.cache.GetCoins()
	alerts, err := h.service.GetAll()
	if err != nil {
		slog.Error("AlertHandler:GetAll", "error", err)
		http.Error(w, "Unable to retrieve alerts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Coins  []model.Coin
		Alerts []model.Alert
	}{
		Coins:  coins,
		Alerts: alerts,
	}

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("AlertHandler:tmpl.Execute", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AlertHandler) Save(w http.ResponseWriter, r *http.Request) {
	var alert model.Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		slog.Error("AlertHandler:Decode", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Save(alert); err != nil {
		slog.Error("AlertHandler:Save", "error", err)
		http.Error(w, "Failed to save alert", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode("Success...")
}

func (h *AlertHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(alerts)
}

func (h *AlertHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid alert ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "Failed to delete alert", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
