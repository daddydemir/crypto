package application

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/daddydemir/crypto/pkg/macd/domain"
)

type MACDApplicationService struct {
	macdRepo  domain.MACDRepository
	priceRepo domain.PriceRepository
}

func NewMACDApplicationService(
	macdRepo domain.MACDRepository,
	priceRepo domain.PriceRepository,
) *MACDApplicationService {
	return &MACDApplicationService{
		macdRepo:  macdRepo,
		priceRepo: priceRepo,
	}
}

func (s *MACDApplicationService) CalculateMACD(symbol string, from, to time.Time) (*domain.MACDResult, error) {
	slog.Info("Calculating MACD", "symbol", symbol, "from", from, "to", to)

	// Get price data
	prices, err := s.priceRepo.GetPriceData(symbol, from, to)
	if err != nil {
		slog.Error("Failed to get price data", "symbol", symbol, "error", err)
		return nil, fmt.Errorf("failed to get price data for %s: %w", symbol, err)
	}

	if len(prices) == 0 {
		return nil, errors.New("no price data available")
	}

	// Create MACD calculator with default parameters
	calculator := domain.NewMACDCalculator()

	// Calculate MACD
	result, err := calculator.Calculate(prices)
	if err != nil {
		slog.Error("Failed to calculate MACD", "symbol", symbol, "error", err)
		return nil, fmt.Errorf("failed to calculate MACD for %s: %w", symbol, err)
	}

	result.Symbol = symbol

	// Save result to repository
	if err := s.macdRepo.SaveMACDResult(symbol, result); err != nil {
		slog.Error("Failed to save MACD result", "symbol", symbol, "error", err)
		// Don't return error, as calculation was successful
	}

	slog.Info("MACD calculation completed",
		"symbol", symbol,
		"dataPoints", len(result.Data),
		"trend", result.Trend,
	)

	return result, nil
}

func (s *MACDApplicationService) CalculateMACDWithCustomParams(
	symbol string,
	from, to time.Time,
	fast, slow, signal int,
) (*domain.MACDResult, error) {
	slog.Info("Calculating MACD with custom params",
		"symbol", symbol,
		"fast", fast,
		"slow", slow,
		"signal", signal,
	)

	// Get price data
	prices, err := s.priceRepo.GetPriceData(symbol, from, to)
	if err != nil {
		slog.Error("Failed to get price data", "symbol", symbol, "error", err)
		return nil, fmt.Errorf("failed to get price data for %s: %w", symbol, err)
	}

	if len(prices) == 0 {
		return nil, errors.New("no price data available")
	}

	// Create MACD calculator with custom parameters
	calculator, err := domain.NewMACDCalculatorWithParams(fast, slow, signal)
	if err != nil {
		slog.Error("Invalid MACD parameters", "error", err)
		return nil, fmt.Errorf("invalid MACD parameters: %w", err)
	}

	// Calculate MACD
	result, err := calculator.Calculate(prices)
	if err != nil {
		slog.Error("Failed to calculate MACD", "symbol", symbol, "error", err)
		return nil, fmt.Errorf("failed to calculate MACD for %s: %w", symbol, err)
	}

	result.Symbol = symbol

	// Save result to repository
	if err := s.macdRepo.SaveMACDResult(symbol, result); err != nil {
		slog.Error("Failed to save MACD result", "symbol", symbol, "error", err)
	}

	return result, nil
}

func (s *MACDApplicationService) GetMACDAnalysis(symbol string) (*domain.MACDAnalysis, error) {
	slog.Info("Getting MACD analysis", "symbol", symbol)

	// Get latest MACD result
	result, err := s.macdRepo.GetLatestMACDResult(symbol)
	if err != nil {
		slog.Error("Failed to get latest MACD result", "symbol", symbol, "error", err)
		return nil, fmt.Errorf("failed to get MACD result for %s: %w", symbol, err)
	}

	if result == nil || len(result.Data) == 0 {
		return nil, errors.New("no MACD data available")
	}

	// Create analysis
	analysis := &domain.MACDAnalysis{
		Symbol:       symbol,
		Result:       result,
		Signal:       result.GetSignal(),
		AnalysisDate: time.Now(),
	}

	// Check for divergences
	if result.IsBullishDivergence() {
		analysis.Divergence = domain.DivergenceBullish
	} else if result.IsBearishDivergence() {
		analysis.Divergence = domain.DivergenceBearish
	} else {
		analysis.Divergence = domain.DivergenceNone
	}

	// Generate recommendation and confidence
	analysis.Recommendation, analysis.Confidence = s.generateRecommendation(result, analysis.Signal, analysis.Divergence)

	slog.Info("MACD analysis completed",
		"symbol", symbol,
		"signal", analysis.Signal,
		"trend", result.Trend,
		"divergence", analysis.Divergence,
		"confidence", analysis.Confidence,
	)

	return analysis, nil
}

func (s *MACDApplicationService) GetTradingSignal(symbol string) (*domain.TradingSignal, error) {
	slog.Info("Getting trading signal", "symbol", symbol)

	// Get latest price data
	prices, err := s.priceRepo.GetLatestPriceData(symbol, 1)
	if err != nil || len(prices) == 0 {
		slog.Error("Failed to get latest price", "symbol", symbol, "error", err)
		return nil, fmt.Errorf("failed to get latest price for %s: %w", symbol, err)
	}

	// Get MACD analysis
	analysis, err := s.GetMACDAnalysis(symbol)
	if err != nil {
		return nil, err
	}

	latestData := analysis.Result.Data[len(analysis.Result.Data)-1]
	currentPrice := prices[0].Price

	// Calculate signal strength
	strength := s.calculateSignalStrength(analysis.Result, analysis.Signal, analysis.Divergence)

	// Generate reason
	reason := s.generateSignalReason(analysis.Signal, analysis.Result.Trend, analysis.Divergence)

	tradingSignal := &domain.TradingSignal{
		Symbol:     symbol,
		Signal:     analysis.Signal,
		Strength:   strength,
		Reason:     reason,
		Price:      currentPrice,
		Timestamp:  time.Now(),
		MACD:       latestData.MACD,
		SignalLine: latestData.Signal,
		Histogram:  latestData.Histogram,
	}

	slog.Info("Trading signal generated",
		"symbol", symbol,
		"signal", tradingSignal.Signal,
		"strength", tradingSignal.Strength,
		"price", tradingSignal.Price,
	)

	return tradingSignal, nil
}

func (s *MACDApplicationService) generateRecommendation(
	result *domain.MACDResult,
	signal domain.Signal,
	divergence domain.DivergenceType,
) (string, float64) {
	var recommendation string
	var confidence float64

	switch signal {
	case domain.SignalBuy:
		if divergence == domain.DivergenceBullish {
			recommendation = "Strong Buy - MACD bullish crossover with bullish divergence detected"
			confidence = 0.85
		} else {
			recommendation = "Buy - MACD bullish crossover detected"
			confidence = 0.70
		}
	case domain.SignalSell:
		if divergence == domain.DivergenceBearish {
			recommendation = "Strong Sell - MACD bearish crossover with bearish divergence detected"
			confidence = 0.85
		} else {
			recommendation = "Sell - MACD bearish crossover detected"
			confidence = 0.70
		}
	case domain.SignalHold:
		if result.Trend == domain.TrendBullish {
			recommendation = "Hold/Watch - Bullish trend but no clear signal"
			confidence = 0.50
		} else if result.Trend == domain.TrendBearish {
			recommendation = "Hold/Watch - Bearish trend but no clear signal"
			confidence = 0.50
		} else {
			recommendation = "Hold - No clear trend or signal"
			confidence = 0.30
		}
	}

	return recommendation, confidence
}

func (s *MACDApplicationService) calculateSignalStrength(
	result *domain.MACDResult,
	signal domain.Signal,
	divergence domain.DivergenceType,
) float64 {
	var strength float64 = 0.5 // Base strength

	// Adjust based on signal type
	switch signal {
	case domain.SignalBuy, domain.SignalSell:
		strength = 0.7
	case domain.SignalHold:
		strength = 0.3
	}

	// Adjust based on divergence
	if divergence != domain.DivergenceNone {
		strength += 0.2
	}

	// Adjust based on trend consistency
	if len(result.Data) >= 3 {
		recent := result.Data[len(result.Data)-3:]
		histogramTrend := s.analyzeHistogramTrend(recent)

		if (signal == domain.SignalBuy && histogramTrend > 0) ||
			(signal == domain.SignalSell && histogramTrend < 0) {
			strength += 0.1
		}
	}

	// Ensure strength is between 0 and 1
	return math.Max(0, math.Min(1, strength))
}

func (s *MACDApplicationService) analyzeHistogramTrend(data []domain.MACDData) float64 {
	if len(data) < 2 {
		return 0
	}

	// Simple trend analysis of histogram
	trend := 0.0
	for i := 1; i < len(data); i++ {
		if data[i].Histogram > data[i-1].Histogram {
			trend += 1
		} else if data[i].Histogram < data[i-1].Histogram {
			trend -= 1
		}
	}

	return trend / float64(len(data)-1)
}

func (s *MACDApplicationService) generateSignalReason(
	signal domain.Signal,
	trend domain.TrendDirection,
	divergence domain.DivergenceType,
) string {
	var reasons []string

	switch signal {
	case domain.SignalBuy:
		reasons = append(reasons, "MACD line crossed above signal line")
	case domain.SignalSell:
		reasons = append(reasons, "MACD line crossed below signal line")
	case domain.SignalHold:
		reasons = append(reasons, "No clear crossover signal")
	}

	if trend != domain.TrendNeutral {
		reasons = append(reasons, fmt.Sprintf("Overall trend is %s", string(trend)))
	}

	if divergence != domain.DivergenceNone {
		reasons = append(reasons, fmt.Sprintf("%s divergence detected", string(divergence)))
	}

	result := ""
	for i, reason := range reasons {
		if i > 0 {
			result += ". "
		}
		result += reason
	}

	return result
}

// CleanupOldData removes old MACD data
func (s *MACDApplicationService) CleanupOldData(olderThan time.Time) error {
	slog.Info("Cleaning up old MACD data", "olderThan", olderThan)

	err := s.macdRepo.DeleteOldMACDData(olderThan)
	if err != nil {
		slog.Error("Failed to cleanup old MACD data", "error", err)
		return fmt.Errorf("failed to cleanup old MACD data: %w", err)
	}

	slog.Info("Old MACD data cleanup completed")
	return nil
}
