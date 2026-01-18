package domain

import "time"

type MACDData struct {
	MACD      float64
	Signal    float64
	Histogram float64
	Date      time.Time
}

type PriceData struct {
	Price float64
	Date  time.Time
}

// Calculate MACD values for given price data
// fastPeriod: typically 12, slowPeriod: typically 26, signalPeriod: typically 9
func Calculate(prices []float64, fastPeriod, slowPeriod, signalPeriod int) []MACDData {
	if len(prices) < slowPeriod {
		return nil
	}

	// Calculate EMAs
	fastEMA := calculateEMA(prices, fastPeriod)
	slowEMA := calculateEMA(prices, slowPeriod)

	// Calculate MACD line (Fast EMA - Slow EMA)
	var macdLine []float64
	minLen := len(fastEMA)
	if len(slowEMA) < minLen {
		minLen = len(slowEMA)
	}

	for i := 0; i < minLen; i++ {
		macdLine = append(macdLine, fastEMA[i]-slowEMA[i])
	}

	// Calculate Signal line (EMA of MACD line)
	signalLine := calculateEMA(macdLine, signalPeriod)

	// Calculate Histogram (MACD - Signal)
	var result []MACDData
	signalStartIndex := len(macdLine) - len(signalLine)

	for i := 0; i < len(signalLine); i++ {
		macdValue := macdLine[signalStartIndex+i]
		signalValue := signalLine[i]
		histogram := macdValue - signalValue

		result = append(result, MACDData{
			MACD:      macdValue,
			Signal:    signalValue,
			Histogram: histogram,
		})
	}

	return result
}

// calculateEMA calculates Exponential Moving Average
func calculateEMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	var ema []float64
	multiplier := 2.0 / float64(period+1)

	// Calculate initial SMA for first EMA value
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema = append(ema, sum/float64(period))

	// Calculate EMA for remaining values
	for i := period; i < len(prices); i++ {
		emaValue := (prices[i] * multiplier) + (ema[len(ema)-1] * (1 - multiplier))
		ema = append(ema, emaValue)
	}

	return ema
}
