package rest

import (
	"encoding/json"
	"errors"
	"github.com/daddydemir/crypto/pkg/purchasegoals/app"
	"github.com/daddydemir/crypto/pkg/purchasegoals/domain"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Handler struct{ app *app.App }

func New(a *app.App) *Handler { return &Handler{app: a} }
func user(r *http.Request) (string, bool) {
	v, ok := r.Context().Value("username").(string)
	return v, ok && v != ""
}
func id(r *http.Request, key string) (uint, error) {
	v, e := strconv.ParseUint(mux.Vars(r)[key], 10, 64)
	return uint(v), e
}
func out(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, e error, status int) {
	if errors.Is(e, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	}
	out(w, status, map[string]string{"message": e.Error()})
}
func decode(r *http.Request, v any) error { return json.NewDecoder(r.Body).Decode(v) }
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	u, ok := user(r)
	if !ok {
		out(w, 401, map[string]string{"message": "unauthorized"})
		return
	}
	v, e := h.app.List(u)
	if e != nil {
		fail(w, e, 500)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "goalId")
	if e != nil {
		fail(w, e, 400)
		return
	}
	v, e := h.app.Get(u, n)
	if e != nil {
		fail(w, e, 500)
		return
	}
	out(w, 200, v)
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	var v domain.Goal
	if e := decode(r, &v); e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.Create(u, v)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 201, o)
}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "goalId")
	var v domain.Goal
	if e == nil {
		e = decode(r, &v)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.Update(u, n, v)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 200, o)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "goalId")
	if e == nil {
		e = h.app.Delete(u, n)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	w.WriteHeader(204)
}
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "goalId")
	var v domain.Item
	if e == nil {
		e = decode(r, &v)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.AddItem(u, n, v)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 201, o)
}
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "itemId")
	var v domain.Item
	if e == nil {
		e = decode(r, &v)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.UpdateItem(u, n, v)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 200, o)
}
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "itemId")
	if e == nil {
		e = h.app.DeleteItem(u, n)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	w.WriteHeader(204)
}
func (h *Handler) Price(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "itemId")
	var v struct {
		Price float64   `json:"price"`
		Date  time.Time `json:"date"`
	}
	if e == nil {
		e = decode(r, &v)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.UpdatePrice(u, n, v.Price, v.Date)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 200, o)
}
func (h *Handler) Saving(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	g, e := id(r, "goalId")
	s := uint(0)
	if mux.Vars(r)["savingId"] != "" && e == nil {
		s, e = id(r, "savingId")
	}
	var v domain.SavingTransaction
	if e == nil {
		e = decode(r, &v)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.Saving(u, g, s, v)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, map[bool]int{true: 201, false: 200}[s == 0], o)
}
func (h *Handler) DeleteSaving(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	g, e := id(r, "goalId")
	var s uint
	if e == nil {
		s, e = id(r, "savingId")
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.DeleteSaving(u, g, s)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 200, o)
}
func (h *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "itemId")
	var v struct {
		ActualPurchasePrice float64   `json:"actualPurchasePrice"`
		PurchaseDate        time.Time `json:"purchaseDate"`
	}
	if e == nil {
		e = decode(r, &v)
	}
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.Purchase(u, n, v.ActualPurchasePrice, v.PurchaseDate)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 200, o)
}
func (h *Handler) CancelItem(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "itemId")
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.CancelItem(u, n)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 200, o)
}
func (h *Handler) CancelGoal(w http.ResponseWriter, r *http.Request) {
	u, _ := user(r)
	n, e := id(r, "goalId")
	if e != nil {
		fail(w, e, 400)
		return
	}
	o, e := h.app.CancelGoal(u, n)
	if e != nil {
		fail(w, e, 400)
		return
	}
	out(w, 200, o)
}
func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/purchase-goals", h.List).Methods("GET")
	r.HandleFunc("/purchase-goals", h.Create).Methods("POST")
	r.HandleFunc("/purchase-goals/{goalId}", h.Get).Methods("GET")
	r.HandleFunc("/purchase-goals/{goalId}", h.Update).Methods("PUT")
	r.HandleFunc("/purchase-goals/{goalId}", h.Delete).Methods("DELETE")
	r.HandleFunc("/purchase-goals/{goalId}/cancel", h.CancelGoal).Methods("POST")
	r.HandleFunc("/purchase-goals/{goalId}/items", h.AddItem).Methods("POST")
	r.HandleFunc("/purchase-items/{itemId}", h.UpdateItem).Methods("PUT")
	r.HandleFunc("/purchase-items/{itemId}", h.DeleteItem).Methods("DELETE")
	r.HandleFunc("/purchase-items/{itemId}/price", h.Price).Methods("POST")
	r.HandleFunc("/purchase-items/{itemId}/purchase", h.Purchase).Methods("POST")
	r.HandleFunc("/purchase-items/{itemId}/cancel", h.CancelItem).Methods("POST")
	r.HandleFunc("/purchase-goals/{goalId}/savings", h.Saving).Methods("POST")
	r.HandleFunc("/purchase-goals/{goalId}/savings/{savingId}", h.Saving).Methods("PUT")
	r.HandleFunc("/purchase-goals/{goalId}/savings/{savingId}", h.DeleteSaving).Methods("DELETE")
}
