package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/daddydemir/crypto/pkg/service"
)

func getExchange(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(service.GetExchange())
	if err != nil {
		slog.Error("getExchange:json.Encode", "error", err)
	}
}

func getExchangeFromDb(w http.ResponseWriter, _ *http.Request) {

	exchangeService := serviceFactory.NewExchangeService()
	response, err := exchangeService.FindAll()
	if err != nil {
		slog.Error("getExchangeFromDb:FindAll", "error", err)
	}
	slog.Info("getExchangeFromDb:service.GetExchangeFromDb", "response", response)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("getExchangeFromDb:json.Encode", "error", err)
	}
}

func createExchange(_ http.ResponseWriter, _ *http.Request) {
	exchangeService := serviceFactory.NewExchangeService()
	exchangeService.CreateExchange()
}
