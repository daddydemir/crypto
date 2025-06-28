package chart

import (
	"github.com/daddydemir/crypto/internal/port/chart"
	"github.com/gorilla/mux"
	"net/http"
)

type BollingerHandler struct {
	serviceFactory func(string) chart.BollingerChart
}

func NewBollingerHandler(factory func(string) chart.BollingerChart) *BollingerHandler {
	return &BollingerHandler{
		serviceFactory: factory,
	}
}

func (h *BollingerHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	coin := vars["coin"]

	service := h.serviceFactory(coin)
	service.CalculateBands()
	service.DrawBands(w, r)
}
