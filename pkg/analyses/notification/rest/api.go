package rest

import (
	"encoding/json"
	"net/http"

	"github.com/daddydemir/crypto/pkg/analyses/notification/app"
	"github.com/gorilla/mux"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	username, _ := r.Context().Value("username").(string)
	notifications := h.app.GetAll(username)
	json.NewEncoder(w).Encode(notifications)
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/notifications", h.GetAll).Methods(http.MethodGet)
}
