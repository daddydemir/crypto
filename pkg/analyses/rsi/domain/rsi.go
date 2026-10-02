package domain

import (
	"math"
	"time"
)

func Calculate(prices []float64) float64 {
	const period = 14
	if len(prices) < period+1 {
		return 0
	}

	start := len(prices) - period - 1
	var gains, losses float64
	for i := start + 1; i < len(prices); i++ {
		diff := prices[i] - prices[i-1]
		if diff >= 0 {
			gains += diff
		} else {
			losses -= diff
		}
	}

	avgGain := gains / period
	avgLoss := losses / period

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return math.Round(rsi*100) / 100
}

type PriceData struct {
	Date  time.Time
	Price float64
}
