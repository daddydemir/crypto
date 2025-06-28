package daily

import (
	"encoding/json"
	"github.com/daddydemir/crypto/internal/domain/dto"
	"github.com/daddydemir/crypto/internal/port/daily"
	"log/slog"
	"net/http"
)

type DailyHandler struct {
	service daily.DailyService
}

func NewDailyHandler(service daily.DailyService) *DailyHandler {
	return &DailyHandler{service: service}
}

func (h *DailyHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	var req dto.DateRangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("GetByDateRange:Decode", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.FindByDateRange(req.StartDate, req.EndDate)
	if err != nil {
		slog.Error("GetByDateRange", "error", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *DailyHandler) GetByIdAndDateRange(w http.ResponseWriter, r *http.Request) {
	var req dto.DateRangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("GetByIdAndDateRange:Decode", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.FindByIdAndDateRange(req.Id, req.StartDate, req.EndDate)
	if err != nil {
		slog.Error("GetByIdAndDateRange", "error", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}
