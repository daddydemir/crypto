package globalsearch

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{ service *Service }

func NewHandler(s *Service) *Handler { return &Handler{s} }
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len([]rune(q)) < 2 {
		json.NewEncoder(w).Encode(Results{Coins: []Item{}, Strategies: []Item{}, Runs: []Item{}, Alerts: []Item{}, Trades: []Item{}})
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	u, _ := r.Context().Value("username").(string)
	v, e := h.service.Search(u, q, limit)
	if e != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"message": e.Error()})
		return
	}
	json.NewEncoder(w).Encode(v)
}
func (h *Handler) Register(r *mux.Router) { r.HandleFunc("/search", h.Search).Methods("GET") }
