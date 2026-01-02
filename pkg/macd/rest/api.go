package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/daddydemir/crypto/pkg/macd/domain"
	"github.com/gorilla/mux"
	"log/slog"
)

type MACDHandler struct {
	service domain.MACDService
}

func NewMACDHandler(service domain.MACDService) *MACDHandler {
	return &MACDHandler{
		service: service,
	}
}

// RegisterRoutes registers all MACD-related routes
func (h *MACDHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/macd/{symbol}/calculate", h.CalculateMACD).Methods("POST")
	router.HandleFunc("/macd/{symbol}/analysis", h.GetMACDAnalysis).Methods("GET")
	router.HandleFunc("/macd/{symbol}/signal", h.GetTradingSignal).Methods("GET")
}

// CalculateMACDRequest represents the request for MACD calculation
type CalculateMACDRequest struct {
	From         string `json:"from"`
	To           string `json:"to"`
	FastPeriod   int    `json:"fast_period,omitempty"`
	SlowPeriod   int    `json:"slow_period,omitempty"`
	SignalPeriod int    `json:"signal_period,omitempty"`
}

// MACDResponse represents the MACD calculation response
type MACDResponse struct {
	Symbol     string                `json:"symbol"`
	Data       []MACDDataResponse    `json:"data"`
	LastMACD   float64               `json:"last_macd"`
	LastSignal float64               `json:"last_signal"`
	Trend      domain.TrendDirection `json:"trend"`
}

type MACDDataResponse struct {
	Date      string  `json:"date"`
	Price     float64 `json:"price"`
	MACD      float64 `json:"macd"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
	EMA12     float64 `json:"ema12"`
	EMA26     float64 `json:"ema26"`
}

// CalculateMACD handles MACD calculation requests
func (h *MACDHandler) CalculateMACD(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	var req CalculateMACDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse dates
	from, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		slog.Error("Invalid from date", "date", req.From, "error", err)
		http.Error(w, "Invalid from date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	to, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		slog.Error("Invalid to date", "date", req.To, "error", err)
		http.Error(w, "Invalid to date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var result *domain.MACDResult

	// Use custom parameters if provided
	if req.FastPeriod > 0 && req.SlowPeriod > 0 && req.SignalPeriod > 0 {
		result, err = h.service.CalculateMACDWithCustomParams(symbol, from, to, req.FastPeriod, req.SlowPeriod, req.SignalPeriod)
	} else {
		result, err = h.service.CalculateMACD(symbol, from, to)
	}

	if err != nil {
		slog.Error("Failed to calculate MACD", "symbol", symbol, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := h.convertMACDResultToResponse(result)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetMACDAnalysis handles MACD analysis requests
func (h *MACDHandler) GetMACDAnalysis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	analysis, err := h.service.GetMACDAnalysis(symbol)
	if err != nil {
		slog.Error("Failed to get MACD analysis", "symbol", symbol, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"symbol":         analysis.Symbol,
		"signal":         analysis.Signal,
		"divergence":     analysis.Divergence,
		"recommendation": analysis.Recommendation,
		"confidence":     analysis.Confidence,
		"analysis_date":  analysis.AnalysisDate.Format(time.RFC3339),
		"macd_result":    h.convertMACDResultToResponse(analysis.Result),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetTradingSignal handles trading signal requests
func (h *MACDHandler) GetTradingSignal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	signal, err := h.service.GetTradingSignal(symbol)
	if err != nil {
		slog.Error("Failed to get trading signal", "symbol", symbol, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"symbol":      signal.Symbol,
		"signal":      signal.Signal,
		"strength":    signal.Strength,
		"reason":      signal.Reason,
		"price":       signal.Price,
		"timestamp":   signal.Timestamp.Format(time.RFC3339),
		"macd":        signal.MACD,
		"signal_line": signal.SignalLine,
		"histogram":   signal.Histogram,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *MACDHandler) convertMACDResultToResponse(result *domain.MACDResult) *MACDResponse {
	response := &MACDResponse{
		Symbol:     result.Symbol,
		Data:       make([]MACDDataResponse, len(result.Data)),
		LastMACD:   result.LastMACD,
		LastSignal: result.LastSignal,
		Trend:      result.Trend,
	}

	for i, data := range result.Data {
		response.Data[i] = MACDDataResponse{
			Date:      data.Date.Format("2006-01-02"),
			Price:     data.Price,
			MACD:      data.MACD,
			Signal:    data.Signal,
			Histogram: data.Histogram,
			EMA12:     data.EMA12,
			EMA26:     data.EMA26,
		}
	}

	return response
}
