package rest

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/strategylab/app"
	"github.com/daddydemir/crypto/pkg/strategylab/domain"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{ a *app.App }

func New(a *app.App) *Handler { return &Handler{a} }
func user(r *http.Request) (string, bool) {
	v, ok := r.Context().Value("username").(string)
	return v, ok && v != ""
}
func id(r *http.Request) (uint, error) {
	v, e := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	return uint(v), e
}
func out(w http.ResponseWriter, s int, v any)   { w.WriteHeader(s); _ = json.NewEncoder(w).Encode(v) }
func bad(w http.ResponseWriter, s int, e error) { out(w, s, map[string]string{"message": e.Error()}) }
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	v, e := h.a.List(u)
	if e != nil {
		bad(w, 500, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e != nil {
		bad(w, 400, e)
		return
	}
	v, e := h.a.Get(u, n)
	if e != nil {
		bad(w, 404, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	var v domain.Strategy
	if e := json.NewDecoder(r.Body).Decode(&v); e != nil {
		bad(w, 400, e)
		return
	}
	v, e := h.a.Create(u, v)
	if e != nil {
		bad(w, 400, e)
		return
	}
	out(w, 201, v)
}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e != nil {
		bad(w, 400, e)
		return
	}
	var v domain.Strategy
	if e = json.NewDecoder(r.Body).Decode(&v); e != nil {
		bad(w, 400, e)
		return
	}
	v, e = h.a.Update(u, n, v)
	if e != nil {
		bad(w, 400, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e == nil {
		e = h.a.Delete(u, n)
	}
	if e != nil {
		bad(w, 404, e)
		return
	}
	w.WriteHeader(204)
}
func (h *Handler) Clone(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e != nil {
		bad(w, 400, e)
		return
	}
	v, e := h.a.Clone(u, n)
	if e != nil {
		bad(w, 400, e)
		return
	}
	out(w, 201, v)
}
func (h *Handler) Start(kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := user(r)
		n, e := id(r)
		if e != nil {
			bad(w, 400, e)
			return
		}
		var q app.RunRequest
		if e = json.NewDecoder(r.Body).Decode(&q); e != nil {
			bad(w, 400, e)
			return
		}
		v, e := h.a.Start(u, n, kind, q)
		if e != nil {
			bad(w, 400, e)
			return
		}
		out(w, 201, v)
	}
}
func (h *Handler) Runs(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	v, e := h.a.Runs(u)
	if e != nil {
		bad(w, 500, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e != nil {
		bad(w, 400, e)
		return
	}
	v, e := h.a.Run(u, n)
	if e != nil {
		bad(w, 404, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) DeleteRun(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e == nil {
		e = h.a.DeleteRun(u, n)
	}
	if e != nil {
		bad(w, 400, e)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) Trades(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e != nil {
		bad(w, 400, e)
		return
	}
	v, e := h.a.Trades(u, n)
	if e != nil {
		bad(w, 500, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Equity(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	if e != nil {
		bad(w, 400, e)
		return
	}
	v, e := h.a.Equity(u, n)
	if e != nil {
		bad(w, 500, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Compare(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	ids := []uint{}
	for _, raw := range strings.Split(r.URL.Query().Get("ids"), ",") {
		v, e := strconv.ParseUint(raw, 10, 64)
		if e != nil {
			bad(w, 400, e)
			return
		}
		ids = append(ids, uint(v))
	}
	v, e := h.a.Compare(u, ids)
	if e != nil {
		bad(w, 400, e)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r)
	status := strings.ToUpper(mux.Vars(r)["action"])
	m := map[string]string{"PAUSE": domain.StatusPaused, "RESUME": domain.StatusRunning, "STOP": domain.StatusStopped}
	status = m[status]
	if status == "" {
		bad(w, 400, assert("unsupported action"))
		return
	}
	if e == nil {
		e = h.a.Status(u, n, status)
	}
	if e != nil {
		bad(w, 400, e)
		return
	}
	w.WriteHeader(204)
}

type assert string

func (a assert) Error() string { return string(a) }
func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/strategies", h.List).Methods("GET")
	r.HandleFunc("/strategies", h.Create).Methods("POST")
	r.HandleFunc("/strategies/{id}", h.Get).Methods("GET")
	r.HandleFunc("/strategies/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/strategies/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/strategies/{id}/clone", h.Clone).Methods("POST")
	r.HandleFunc("/strategies/{id}/backtests", h.Start(domain.RunBacktest)).Methods("POST")
	r.HandleFunc("/strategies/{id}/forward-tests", h.Start(domain.RunForward)).Methods("POST")
	r.HandleFunc("/strategy-runs", h.Runs).Methods("GET")
	r.HandleFunc("/strategy-runs/compare", h.Compare).Methods("GET")
	r.HandleFunc("/strategy-runs/{id}", h.Run).Methods("GET")
	r.HandleFunc("/strategy-runs/{id}", h.DeleteRun).Methods("DELETE")
	r.HandleFunc("/strategy-runs/{id}/trades", h.Trades).Methods("GET")
	r.HandleFunc("/strategy-runs/{id}/equity", h.Equity).Methods("GET")
	r.HandleFunc("/strategy-runs/{id}/{action:pause|resume|stop}", h.Status).Methods("POST")
}
