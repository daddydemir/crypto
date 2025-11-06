package indicator

import "testing"

func TestRSI_Calculate(t *testing.T) {
	rsi := NewRSI()
	prices := []float64{132.04, 130.59, 131.34, 137.07, 141.30, 143.21, 142.7, 138.89, 139.79, 144.65, 145.37, 144.49, 154.48, 164.62}

	result := rsi.Calculate(prices)

	if result != 85.51 {
		t.Errorf("actual: %v, expected: %v", result, "85.51")
	}
}
