package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/atr/app"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) *Handler {
	return &Handler{
		app: app,
	}
}

func (h *Handler) Points(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
	}

	points, err := h.app.GetPoints(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if len(points) == 0 {
		http.Error(w, "no points found", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(points)
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/adi/coin/{symbol}", h.Points)
}
