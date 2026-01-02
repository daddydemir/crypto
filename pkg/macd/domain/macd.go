package domain

import (
	"errors"
	"time"
)

// MACDData represents a single MACD calculation result
type MACDData struct {
	Date      time.Time
	Price     float64
	MACD      float64
	Signal    float64
	Histogram float64
	EMA12     float64
	EMA26     float64
}

// MACDCalculator contains the parameters for MACD calculation
type MACDCalculator struct {
	FastPeriod   int // Usually 12
	SlowPeriod   int // Usually 26
	SignalPeriod int // Usually 9
}

// MACDResult represents the complete MACD analysis
type MACDResult struct {
	Symbol     string
	Data       []MACDData
	LastMACD   float64
	LastSignal float64
	Trend      TrendDirection
}

type TrendDirection string

const (
	TrendBullish TrendDirection = "BULLISH"
	TrendBearish TrendDirection = "BEARISH"
	TrendNeutral TrendDirection = "NEUTRAL"
)

// NewMACDCalculator creates a new MACD calculator with default parameters
func NewMACDCalculator() *MACDCalculator {
	return &MACDCalculator{
		FastPeriod:   12,
		SlowPeriod:   26,
		SignalPeriod: 9,
	}
}

// NewMACDCalculatorWithParams creates a MACD calculator with custom parameters
func NewMACDCalculatorWithParams(fast, slow, signal int) (*MACDCalculator, error) {
	if fast >= slow {
		return nil, errors.New("fast period must be less than slow period")
	}
	if fast <= 0 || slow <= 0 || signal <= 0 {
		return nil, errors.New("all periods must be positive")
	}

	return &MACDCalculator{
		FastPeriod:   fast,
		SlowPeriod:   slow,
		SignalPeriod: signal,
	}, nil
}

// Calculate performs MACD calculation on price data
func (mc *MACDCalculator) Calculate(prices []PriceData) (*MACDResult, error) {
	if len(prices) < mc.SlowPeriod+mc.SignalPeriod {
		return nil, errors.New("insufficient data for MACD calculation")
	}

	// Calculate EMAs
	ema12 := mc.calculateEMA(prices, mc.FastPeriod)
	ema26 := mc.calculateEMA(prices, mc.SlowPeriod)

	// Calculate MACD line
	var macdData []MACDData
	var macdValues []float64

	for i := mc.SlowPeriod - 1; i < len(prices); i++ {
		macd := ema12[i] - ema26[i]
		macdValues = append(macdValues, macd)

		data := MACDData{
			Date:  prices[i].Date,
			Price: prices[i].Price,
			MACD:  macd,
			EMA12: ema12[i],
			EMA26: ema26[i],
		}
		macdData = append(macdData, data)
	}

	// Calculate Signal line (EMA of MACD)
	signalEMA := mc.calculateEMAFromValues(macdValues, mc.SignalPeriod)

	// Update MACD data with signal and histogram
	for i := mc.SignalPeriod - 1; i < len(macdData); i++ {
		macdData[i].Signal = signalEMA[i]
		macdData[i].Histogram = macdData[i].MACD - macdData[i].Signal
	}

	// Determine trend
	trend := mc.determineTrend(macdData)

	result := &MACDResult{
		Data:  macdData[mc.SignalPeriod-1:], // Only return data with complete calculations
		Trend: trend,
	}

	if len(result.Data) > 0 {
		lastData := result.Data[len(result.Data)-1]
		result.LastMACD = lastData.MACD
		result.LastSignal = lastData.Signal
	}

	return result, nil
}

// calculateEMA calculates Exponential Moving Average
func (mc *MACDCalculator) calculateEMA(prices []PriceData, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	ema := make([]float64, len(prices))
	multiplier := 2.0 / (float64(period) + 1.0)

	// Calculate initial SMA for the first EMA value
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i].Price
	}
	ema[period-1] = sum / float64(period)

	// Calculate EMA for the rest
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i].Price * multiplier) + (ema[i-1] * (1 - multiplier))
	}

	return ema
}

// calculateEMAFromValues calculates EMA from a slice of values
func (mc *MACDCalculator) calculateEMAFromValues(values []float64, period int) []float64 {
	if len(values) < period {
		return nil
	}

	ema := make([]float64, len(values))
	multiplier := 2.0 / (float64(period) + 1.0)

	// Calculate initial SMA for the first EMA value
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += values[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate EMA for the rest
	for i := period; i < len(values); i++ {
		ema[i] = (values[i] * multiplier) + (ema[i-1] * (1 - multiplier))
	}

	return ema
}

// determineTrend analyzes the MACD data to determine the current trend
func (mc *MACDCalculator) determineTrend(data []MACDData) TrendDirection {
	if len(data) < 3 {
		return TrendNeutral
	}

	recent := data[len(data)-3:]

	// Check for bullish signals
	if recent[2].MACD > recent[2].Signal && recent[2].Histogram > recent[1].Histogram {
		return TrendBullish
	}

	// Check for bearish signals
	if recent[2].MACD < recent[2].Signal && recent[2].Histogram < recent[1].Histogram {
		return TrendBearish
	}

	return TrendNeutral
}

// GetSignal returns trading signal based on MACD
func (mr *MACDResult) GetSignal() Signal {
	if len(mr.Data) < 2 {
		return SignalHold
	}

	current := mr.Data[len(mr.Data)-1]
	previous := mr.Data[len(mr.Data)-2]

	// MACD line crosses above signal line
	if previous.MACD <= previous.Signal && current.MACD > current.Signal {
		return SignalBuy
	}

	// MACD line crosses below signal line
	if previous.MACD >= previous.Signal && current.MACD < current.Signal {
		return SignalSell
	}

	return SignalHold
}

// IsBullishDivergence checks for bullish divergence
func (mr *MACDResult) IsBullishDivergence() bool {
	if len(mr.Data) < 10 {
		return false
	}

	// Simple implementation - can be enhanced
	recent := mr.Data[len(mr.Data)-5:]

	// Price making lower lows while MACD making higher lows
	priceDecreasing := recent[4].Price < recent[0].Price
	macdIncreasing := recent[4].MACD > recent[0].MACD

	return priceDecreasing && macdIncreasing
}

// IsBearishDivergence checks for bearish divergence
func (mr *MACDResult) IsBearishDivergence() bool {
	if len(mr.Data) < 10 {
		return false
	}

	// Simple implementation - can be enhanced
	recent := mr.Data[len(mr.Data)-5:]

	// Price making higher highs while MACD making lower highs
	priceIncreasing := recent[4].Price > recent[0].Price
	macdDecreasing := recent[4].MACD < recent[0].MACD

	return priceIncreasing && macdDecreasing
}

// PriceData represents a single price point
type PriceData struct {
	Date  time.Time
	Price float64
}

// Signal represents trading signals
type Signal string

const (
	SignalBuy  Signal = "BUY"
	SignalSell Signal = "SELL"
	SignalHold Signal = "HOLD"
)
