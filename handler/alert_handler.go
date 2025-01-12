package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/assets"
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/daddydemir/crypto/pkg/service"
	"io"
	"log/slog"
	"net/http"
)

var alertService *service.AlertService
var cacheService *service.CacheService

func alertPage(w http.ResponseWriter, r *http.Request) {
	tmpl := assets.GetTemplate("templates/alert.html")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	coins := cacheService.GetCoins()

	alerts, err := alertService.GetAll()
	if err != nil {
		slog.Error("alertService.GetAll", "error", err)
		http.Error(w, "Unable to retrieve alerts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Coins  []coincap.Coin
		Alerts []model.Alert
	}{
		Coins:  coins,
		Alerts: alerts,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error("tmpl.Execute", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func alert(w http.ResponseWriter, r *http.Request) {
	all, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("ReadAll", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req model.Alert
	err = json.Unmarshal(all, &req)
	if err != nil {
		slog.Error("Unmarshall", "error", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err = alertService.Save(req)
	if err != nil {
		slog.Error("Save", "error", err)
		http.Error(w, "Failed to save alert", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode("Success...")
	if err != nil {
		slog.Error("Encode", "error", err)
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}
