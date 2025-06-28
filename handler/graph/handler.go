package graph

import (
	"github.com/daddydemir/crypto/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type GraphHandler struct {
	rsiFactory func(coin string) service.RSIService
	emaFactory func(coin string) service.ChartDrawer
	maFactory  func(coin string) service.MovingAverageService
}

func NewGraphHandler(
	rsiFactory func(coin string) service.RSIService,
	emaFactory func(coin string) service.ChartDrawer,
	maFactory func(coin string) service.MovingAverageService,
) *GraphHandler {
	return &GraphHandler{
		rsiFactory: rsiFactory,
		emaFactory: emaFactory,
		maFactory:  maFactory,
	}
}
func (h *GraphHandler) RSIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coin := vars["coin"]

	service := h.rsiFactory(coin) // soyut servis
	handler := service.Draw()

	if handler == nil {
		http.Error(w, "RSI graph cannot be drawn", http.StatusInternalServerError)
		return
	}

	handler(w, r)
}

func (h *GraphHandler) EMAHandler(w http.ResponseWriter, r *http.Request) {
	coin := mux.Vars(r)["coin"]
	service := h.emaFactory(coin)
	handler := service.Draw()

	if handler == nil {
		http.Error(w, "EMA graph cannot be drawn", http.StatusInternalServerError)
		return
	}

	handler(w, r)
}
