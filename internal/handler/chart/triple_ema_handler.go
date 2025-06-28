package chart

import (
	"net/http"

	"github.com/daddydemir/crypto/internal/port/chart"
	"github.com/gorilla/mux"
)

type TripleEmaHandler struct {
	serviceFactory func(coin string) chart.TripleEmaChart
}

func NewTripleEmaHandler(factory func(coin string) chart.TripleEmaChart) *TripleEmaHandler {
	return &TripleEmaHandler{
		serviceFactory: factory,
	}
}

func (h *TripleEmaHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	coin := vars["coin"]

	service := h.serviceFactory(coin)
	service.Draw(w, r)
}
