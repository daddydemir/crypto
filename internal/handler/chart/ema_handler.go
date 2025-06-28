package chart

import (
	"github.com/daddydemir/crypto/internal/port/chart"
	"github.com/gorilla/mux"
	"net/http"
)

type EmaHandler struct {
	serviceFactory func(coin string) chart.EmaChart
}

func NewEmaHandler(factory func(coin string) chart.EmaChart) *EmaHandler {
	return &EmaHandler{
		serviceFactory: factory,
	}
}

func (h *EmaHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	coin := vars["coin"]

	service := h.serviceFactory(coin)
	service.Draw(w, r)
}
