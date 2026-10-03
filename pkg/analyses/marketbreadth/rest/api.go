package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	httpcfg "github.com/daddydemir/crypto/config/http"
	"github.com/daddydemir/crypto/pkg/analyses/marketbreadth/app"
	"github.com/daddydemir/crypto/pkg/analyses/marketbreadth/domain"
	"github.com/gorilla/mux"
)

type Handler struct{ app *app.App }

func NewHandler(a *app.App) *Handler { return &Handler{app: a} }
func timeframe(r *http.Request) (string, bool) {
	v := strings.ToLower(r.URL.Query().Get("timeframe"))
	if v == "" {
		v = "1d"
	}
	return v, v == "1d"
}
func fail(w http.ResponseWriter, status int, message string) {
	httpcfg.WriteJSONError(w, status, httpcfg.NewHttpError(message))
}
func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	_, ok := timeframe(r)
	if !ok {
		fail(w, 400, "only the 1d timeframe is supported")
		return
	}
	result, err := h.app.Current()
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	json.NewEncoder(w).Encode(result)
}
func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	_, ok := timeframe(r)
	if !ok {
		fail(w, 400, "only the 1d timeframe is supported")
		return
	}
	period := strings.ToLower(r.URL.Query().Get("period"))
	if period == "" {
		period = "30d"
	}
	if period != "30d" && period != "90d" && period != "1y" {
		fail(w, 400, "period must be one of 30d, 90d, 1y")
		return
	}
	result, err := h.app.History(period)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	json.NewEncoder(w).Encode(result)
}
func (h *Handler) Coins(w http.ResponseWriter, r *http.Request) {
	_, ok := timeframe(r)
	if !ok {
		fail(w, 400, "only the 1d timeframe is supported")
		return
	}
	metric := domain.Metric(strings.ToLower(r.URL.Query().Get("metric")))
	if !domain.ValidMetrics[metric] {
		fail(w, 400, "unsupported metric")
		return
	}
	result, err := h.app.Coins(metric)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	json.NewEncoder(w).Encode(result)
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/market-breadth", h.Current).Methods(http.MethodGet)
	router.HandleFunc("/market-breadth/history", h.History).Methods(http.MethodGet)
	router.HandleFunc("/market-breadth/coins", h.Coins).Methods(http.MethodGet)
}
