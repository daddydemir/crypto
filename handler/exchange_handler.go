package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/database/service"
	"log"
	"net/http"
)

func getExchange(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(service.GetExchange())
	if err != nil {
		log.Println("::getExchange:: err:{}", err)
	}
}

func getExchangeFromDb(w http.ResponseWriter, r *http.Request) {
	response := service.GetExchangeFromDb()
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("::getExchangeFromDb:: err:{}", err)
	}
}

func createExchange(w http.ResponseWriter, r *http.Request) {
	service.CreateExchange()
}
