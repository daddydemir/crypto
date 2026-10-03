package tradeexplorer

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct{ app *App }

func NewHandler(app *App) *Handler    { return &Handler{app} }
func username(r *http.Request) string { v, _ := r.Context().Value("username").(string); return v }
func write(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func parseDate(v string) *time.Time {
	if v == "" {
		return nil
	}
	d, e := time.Parse("2006-01-02", v)
	if e != nil {
		return nil
	}
	return &d
}
func filters(r *http.Request) Filters {
	q := r.URL.Query()
	sid, _ := strconv.ParseUint(q.Get("strategyId"), 10, 64)
	rid, _ := strconv.ParseUint(q.Get("runId"), 10, 64)
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("pageSize"))
	return Filters{Result: strings.ToLower(q.Get("result")), Coin: q.Get("coin"), RunType: q.Get("runType"), StrategyID: uint(sid), RunID: uint(rid), From: parseDate(q.Get("from")), To: parseDate(q.Get("to")), Sort: q.Get("sort"), Desc: q.Get("direction") != "asc", Page: page, PageSize: size}
}
func tradeID(r *http.Request) (uint, error) {
	v, e := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	return uint(v), e
}
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	v, e := h.app.List(username(r), filters(r))
	if e != nil {
		write(w, 500, map[string]string{"message": e.Error()})
		return
	}
	write(w, 200, v)
}
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	id, e := tradeID(r)
	if e != nil {
		write(w, 400, map[string]string{"message": e.Error()})
		return
	}
	v, e := h.app.Detail(username(r), id)
	if e != nil {
		write(w, 404, map[string]string{"message": e.Error()})
		return
	}
	write(w, 200, v)
}
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	v, e := h.app.Statistics(username(r), filters(r))
	if e != nil {
		write(w, 500, map[string]string{"message": e.Error()})
		return
	}
	write(w, 200, v)
}
func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	v, e := h.app.Options(username(r))
	if e != nil {
		write(w, 500, map[string]string{"message": e.Error()})
		return
	}
	write(w, 200, v)
}
func (h *Handler) Note(w http.ResponseWriter, r *http.Request) {
	id, e := tradeID(r)
	var body struct {
		Note string `json:"note"`
	}
	if e == nil {
		e = json.NewDecoder(r.Body).Decode(&body)
	}
	if e == nil {
		e = h.app.UpdateNote(username(r), id, body.Note)
	}
	if e != nil {
		write(w, 400, map[string]string{"message": e.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/trade-explorer/trades", h.List).Methods("GET")
	r.HandleFunc("/trade-explorer/trades/{id}", h.Detail).Methods("GET")
	r.HandleFunc("/trade-explorer/trades/{id}/note", h.Note).Methods("PATCH")
	r.HandleFunc("/trade-explorer/statistics", h.Stats).Methods("GET")
	r.HandleFunc("/trade-explorer/options", h.Options).Methods("GET")
}
