package chart

import (
	"net/http"

	"github.com/daddydemir/crypto/internal/port/chart"
	"github.com/gorilla/mux"
)

type RsiHandler struct {
	serviceFactory func(coin string) chart.RsiChart
}

func NewRsiHandler(factory func(coin string) chart.RsiChart) *RsiHandler {
	return &RsiHandler{
		serviceFactory: factory,
	}
}

func (h *RsiHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	coin := vars["coin"]

	service := h.serviceFactory(coin)
	service.Draw(w, r)
}
