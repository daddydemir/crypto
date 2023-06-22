package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/database/service"
	"log"
	"net/http"
)

func dailyStart(w http.ResponseWriter, r *http.Request) {
	service.CreateDaily(true)
}

func dailyEnd(w http.ResponseWriter, r *http.Request) {
	service.CreateDaily(false)
}

func daily(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(service.GetDaily())
	if err != nil {
		log.Println("::daily:: err:{}", err)
	}
}
