package chart

import (
	"net/http"

	"github.com/daddydemir/crypto/internal/port/chart"
	"github.com/gorilla/mux"
)

type SmaHandler struct {
	serviceFactory func(coin string) chart.SmaChart
}

func NewSmaHandler(factory func(coin string) chart.SmaChart) *SmaHandler {
	return &SmaHandler{
		serviceFactory: factory,
	}
}

func (h *SmaHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	coin := vars["coin"]

	service := h.serviceFactory(coin)
	service.Draw(w, r)
}
