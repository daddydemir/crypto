package handler

import (
	"github.com/daddydemir/crypto/pkg/database/service"
	"net/http"
)

func getWeekly(_ http.ResponseWriter, _ *http.Request) {
	service.CreateWeekly()
}
