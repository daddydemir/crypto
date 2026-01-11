package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/analyses/adi/app"
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

func (h *Handler) GetADI(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}

	series, err := h.app.GetADISeries(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(series) == 0 {
		http.Error(w, "no data found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(series)
}

// RegisterRoutes registers ADI routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/adi/coin/{symbol}", h.GetADI)
}
