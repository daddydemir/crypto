package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/atr/application"
	"github.com/gorilla/mux"
	"net/http"
)

type AtrHandler struct {
	service *application.PointService
}

func NewAtrHandler(service *application.PointService) *AtrHandler {
	return &AtrHandler{
		service: service,
	}
}

func (handler *AtrHandler) Points(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
	}

	points, err := handler.service.GetPoints(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if len(points) == 0 {
		http.Error(w, "no points found", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(points)
	}
}
