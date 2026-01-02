package domain

import "time"

// MACDRepository defines the contract for MACD data persistence
type MACDRepository interface {
	SaveMACDResult(symbol string, result *MACDResult) error
	GetMACDResult(symbol string, from, to time.Time) (*MACDResult, error)
	GetLatestMACDResult(symbol string) (*MACDResult, error)
	DeleteOldMACDData(olderThan time.Time) error
}

// PriceRepository defines the contract for price data retrieval
type PriceRepository interface {
	GetPriceData(symbol string, from, to time.Time) ([]PriceData, error)
	GetLatestPriceData(symbol string, limit int) ([]PriceData, error)
}

// MACDService defines the business logic interface
type MACDService interface {
	CalculateMACD(symbol string, from, to time.Time) (*MACDResult, error)
	CalculateMACDWithCustomParams(symbol string, from, to time.Time, fast, slow, signal int) (*MACDResult, error)
	GetMACDAnalysis(symbol string) (*MACDAnalysis, error)
	GetTradingSignal(symbol string) (*TradingSignal, error)
}

// MACDAnalysis represents a complete MACD analysis with additional insights
type MACDAnalysis struct {
	Symbol         string
	Result         *MACDResult
	Signal         Signal
	Divergence     DivergenceType
	Recommendation string
	Confidence     float64
	AnalysisDate   time.Time
}

// TradingSignal represents a trading signal with context
type TradingSignal struct {
	Symbol     string
	Signal     Signal
	Strength   float64
	Reason     string
	Price      float64
	Timestamp  time.Time
	MACD       float64
	SignalLine float64
	Histogram  float64
}

// DivergenceType represents the type of divergence
type DivergenceType string

const (
	DivergenceBullish DivergenceType = "BULLISH"
	DivergenceBearish DivergenceType = "BEARISH"
	DivergenceNone    DivergenceType = "NONE"
)
