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
		return
	}

	points, err := h.app.GetPoints(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(points) == 0 {
		http.Error(w, "no points found", http.StatusNotFound)
		return
	} else {
		json.NewEncoder(w).Encode(points)
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/atr/coin/{symbol}", h.Points)
}
