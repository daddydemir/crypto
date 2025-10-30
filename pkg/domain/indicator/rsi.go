package indicator

import (
	"math"
	"time"
)

type RSI struct{}

func NewRSI() *RSI {
	return &RSI{}
}

func (_ *RSI) Calculate(prices []float64) float64 {
	if len(prices) < 14 {
		return 0
	}

	var gains, losses float64
	for i := 1; i < len(prices); i++ {
		diff := prices[i] - prices[i-1]
		if diff >= 0 {
			gains += diff
		} else {
			losses -= diff
		}
	}

	avgGain := gains / 14
	avgLoss := losses / 14

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
