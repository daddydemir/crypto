package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	httpResp "github.com/daddydemir/crypto/config/http"
	"github.com/daddydemir/crypto/pkg/auth/basic/app"
	"github.com/daddydemir/crypto/pkg/auth/basic/domain"
	"github.com/gorilla/mux"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("CreateUser:json.NewDecoder", "err", err)
		httpResp.WriteJSONError(w, http.StatusBadRequest, httpResp.NewHttpError(err.Error()))
		return
	}

	err = h.app.CreateUser(user)
	if err != nil {
		slog.Error("CreateUser:h.app.CreateUser", "err", err)
		httpResp.WriteJSONError(w, http.StatusInternalServerError, httpResp.NewHttpError(err.Error()))
		return
	}
	httpResp.WriteJSONError(w, http.StatusOK, httpResp.NewHttpError("user created"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("Login:json.NewDecoder", "err", err)
		httpResp.WriteJSONError(w, http.StatusBadRequest, httpResp.NewHttpError(err.Error()))
		return
	}

	response, err := h.app.Login(user.Username, user.Password)
	if err != nil {
		slog.Error("Login:h.app.Login", "err", err)
		httpResp.WriteJSONError(w, http.StatusInternalServerError, httpResp.NewHttpError(err.Error()))
		return
	}
	httpResp.WriteJSONError(w, http.StatusOK, response)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
}
