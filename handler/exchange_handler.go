package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/pkg/database/service"
	"net/http"
)

func getExchange(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(service.GetExchange())
	if err != nil {
		log.Errorln("::getExchange:: err:{}", err)
	}
}

func getExchangeFromDb(w http.ResponseWriter, r *http.Request) {
	response := service.GetExchangeFromDb()
	log.Infoln("::getExchangeFromDb:: response:{}", response)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln("::getExchangeFromDb:: err:{}", err)
	}
}

func createExchange(w http.ResponseWriter, r *http.Request) {
	service.CreateExchange()
}
