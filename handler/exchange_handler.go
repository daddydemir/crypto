package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/database/service"
	"log/slog"
	"net/http"
)

func getExchange(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(service.GetExchange())
	if err != nil {
		slog.Error("getExchange:json.Encode", "error", err)
	}
}

func getExchangeFromDb(w http.ResponseWriter, _ *http.Request) {
	response := service.GetExchangeFromDb()
	slog.Info("getExchangeFromDb:service.GetExchangeFromDb", "response", response)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("getExchangeFromDb:json.Encode", "error", err)
	}
}

func createExchange(_ http.ResponseWriter, _ *http.Request) {
	service.CreateExchange()
}
