package handler

import (
	"github.com/daddydemir/crypto/pkg/database/service"
	"net/http"
)

func getWeekly(w http.ResponseWriter, r *http.Request) {
	service.CreateWeekly()
}
