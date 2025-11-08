package exponentialma

import (
	"testing"
	"time"
)

func TestCalculateSeries(t *testing.T) {
	baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	dates := make([]time.Time, 100)
	for i := range dates {
		dates[i] = baseDate.AddDate(0, 0, i)
	}

	t.Run("empty inputs", func(t *testing.T) {
		result := CalculateSeries(nil, nil)
		if result != nil {
			t.Errorf("Expected nil for empty inputs, got %v", result)
		}
	})

	t.Run("mismatched lengths", func(t *testing.T) {
		prices := []float64{1.0, 2.0}
		wrongDates := []time.Time{baseDate}

		result := CalculateSeries(wrongDates, prices)
		if result != nil {
			t.Errorf("Expected nil for mismatched lengths, got %v", result)
		}
	})

	t.Run("insufficient data for MA99", func(t *testing.T) {
		shortDates := make([]time.Time, 50)
		shortPrices := make([]float64, 50)
		for i := range shortDates {
			shortDates[i] = baseDate.AddDate(0, 0, i)
			shortPrices[i] = float64(i + 1)
		}

		result := CalculateSeries(shortDates, shortPrices)
		if result != nil {
			t.Errorf("Expected nil for insufficient data, got %v", result)
		}
	})

	t.Run("valid data with increasing prices", func(t *testing.T) {
		prices := make([]float64, 100)
		for i := range prices {
			prices[i] = float64(i + 1)
		}

		result := CalculateSeries(dates, prices)

		if len(result) != 1 {
			t.Errorf("Expected 1 result, got %d", len(result))
			return
		}

		point := result[0]
		if point.MA7 <= point.MA25 || point.MA25 <= point.MA99 {
			t.Errorf("Expected MA7 > MA25 > MA99, got MA7=%.2f, MA25=%.2f, MA99=%.2f",
				point.MA7, point.MA25, point.MA99)
		}

		lastPrice := prices[len(prices)-1]
		if point.MA7 > lastPrice || point.MA25 > lastPrice || point.MA99 > lastPrice {
			t.Errorf("All MAs should be <= last price (%.2f), got MA7=%.2f, MA25=%.2f, MA99=%.2f",
				lastPrice, point.MA7, point.MA25, point.MA99)
		}
	})

	t.Run("multiple result points", func(t *testing.T) {
		extendedDates := make([]time.Time, 150)
		extendedPrices := make([]float64, 150)
		for i := range extendedDates {
			extendedDates[i] = baseDate.AddDate(0, 0, i)
			extendedPrices[i] = float64(i + 1)
		}

		result := CalculateSeries(extendedDates, extendedPrices)
		expectedPoints := len(extendedPrices) - 99

		if len(result) != expectedPoints {
			t.Errorf("Expected %d result points, got %d", expectedPoints, len(result))
		}

		for i, point := range result {
			expectedDate := extendedDates[99+i]
			if !point.Date.Equal(expectedDate) {
				t.Errorf("Point %d: expected date %v, got %v", i, expectedDate, point.Date)
			}
		}
	})
}

func TestEMA(t *testing.T) {
	t.Run("single element", func(t *testing.T) {
		prices := []float64{10.0}
		period := 7

		result := ema(prices, period)

		if len(result) != 1 {
			t.Errorf("Expected 1 result, got %d", len(result))
		}
		if result[0] != prices[0] {
			t.Errorf("Expected first element to be %.2f, got %.2f", prices[0], result[0])
		}
	})

	t.Run("constant values", func(t *testing.T) {
		prices := []float64{5.0, 5.0, 5.0, 5.0, 5.0}
		period := 3

		result := ema(prices, period)

		for i, val := range result {
			if val != 5.0 {
				t.Errorf("Expected all values to be 5.0, got result[%d]=%.2f", i, val)
			}
		}
	})

	t.Run("increasing values", func(t *testing.T) {
		prices := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		period := 3

		result := ema(prices, period)

		// EMA'nın artan olması gerekir
		for i := 1; i < len(result); i++ {
			if result[i] <= result[i-1] {
				t.Errorf("EMA should be increasing, got result[%d]=%.2f <= result[%d]=%.2f",
					i, result[i], i-1, result[i-1])
			}
		}

		// EMA son değerden küçük olmalı
		lastPrice := prices[len(prices)-1]
		if result[len(result)-1] >= lastPrice {
			t.Errorf("Last EMA (%.2f) should be less than last price (%.2f)",
				result[len(result)-1], lastPrice)
		}
	})

	t.Run("different periods", func(t *testing.T) {
		prices := []float64{10.0, 12.0, 11.0, 13.0, 14.0, 15.0, 13.0, 12.0, 11.0, 10.0}

		// Daha kısa periyot daha hızlı tepki vermeli
		ema3 := ema(prices, 3)
		ema10 := ema(prices, 10)

		// EMA3'ün EMA10'dan daha değişken olması gerekir
		if ema3[len(ema3)-1] == ema10[len(ema10)-1] {
			t.Error("Different period EMAs should have different values")
		}
	})
}

func BenchmarkCalculateSeries(b *testing.B) {
	// 1000 günlük test verisi
	dates := make([]time.Time, 1000)
	prices := make([]float64, 1000)
	baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := range dates {
		dates[i] = baseDate.AddDate(0, 0, i)
		prices[i] = float64(i%100 + 1) // 1-100 arası dalgalanan fiyatlar
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateSeries(dates, prices)
	}
}

func BenchmarkEMA(b *testing.B) {
	prices := make([]float64, 1000)
	for i := range prices {
		prices[i] = float64(i%100 + 1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ema(prices, 25)
	}
}
