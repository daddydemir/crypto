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
		log.Println("::daily:: err:{}", err)
	}
}
