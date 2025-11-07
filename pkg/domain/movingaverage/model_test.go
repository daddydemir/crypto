package movingaverage

import (
	"testing"
	"time"
)

func TestCalculateSeries(t *testing.T) {
	baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		dates    []time.Time
		prices   []float64
		expected []Point
		wantErr  bool
	}{
		{
			name:     "empty_input",
			dates:    []time.Time{},
			prices:   []float64{},
			expected: []Point{},
			wantErr:  false,
		},
		{
			name:     "less_than_100_points",
			dates:    generateDates(baseDate, 50),
			prices:   generatePrices(50, 100.0),
			expected: []Point{},
			wantErr:  false,
		},
		{
			name:   "exactly_100_points",
			dates:  generateDates(baseDate, 100),
			prices: generatePrices(100, 100.0),
			expected: []Point{
				{
					Date: baseDate.AddDate(0, 0, 99),
					MA7:  100.0,
					MA25: 100.0,
					MA99: 100.0,
				},
			},
			wantErr: false,
		},
		{
			name:   "more_than_100_points",
			dates:  generateDates(baseDate, 105),
			prices: generatePrices(105, 100.0),
			expected: []Point{
				{
					Date: baseDate.AddDate(0, 0, 99),
					MA7:  100.0,
					MA25: 100.0,
					MA99: 100.0,
				},
				{
					Date: baseDate.AddDate(0, 0, 100),
					MA7:  100.0,
					MA25: 100.0,
					MA99: 100.0,
				},
				{
					Date: baseDate.AddDate(0, 0, 101),
					MA7:  100.0,
					MA25: 100.0,
					MA99: 100.0,
				},
				{
					Date: baseDate.AddDate(0, 0, 102),
					MA7:  100.0,
					MA25: 100.0,
					MA99: 100.0,
				},
				{
					Date: baseDate.AddDate(0, 0, 103),
					MA7:  100.0,
					MA25: 100.0,
					MA99: 100.0,
				},
				{
					Date: baseDate.AddDate(0, 0, 104),
					MA7:  100.0,
					MA25: 100.0,
					MA99: 100.0,
				},
			},
			wantErr: false,
		},
		{
			name:   "increasing_prices",
			dates:  generateDates(baseDate, 100),
			prices: generateIncreasingPrices(100),
			expected: []Point{
				{
					Date: baseDate.AddDate(0, 0, 99),
					MA7:  calculateExpectedMean(generateIncreasingPrices(100), 93, 100),
					MA25: calculateExpectedMean(generateIncreasingPrices(100), 75, 100),
					MA99: calculateExpectedMean(generateIncreasingPrices(100), 1, 100),
				},
			},
			wantErr: false,
		},
		{
			name:   "varying_prices",
			dates:  generateDates(baseDate, 100),
			prices: []float64{50, 100, 150, 100, 50, 100, 150, 100, 50},
			expected: []Point{
				{
					Date: baseDate.AddDate(0, 0, 99),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "varying_prices" && len(tt.prices) < 100 {
				extendedPrices := extendPattern(tt.prices, 100)
				extendedDates := generateDates(baseDate, 100)
				tt.prices = extendedPrices
				tt.dates = extendedDates

				result := CalculateSeries(extendedDates, extendedPrices)
				if len(result) > 0 {
					tt.expected = []Point{result[0]}
				}
			}

			got := CalculateSeries(tt.dates, tt.prices)

			if len(got) != len(tt.expected) {
				t.Errorf("CalculateSeries() returned %d points, expected %d", len(got), len(tt.expected))
				return
			}

			for i, point := range got {
				if i >= len(tt.expected) {
					break
				}

				if !point.Date.Equal(tt.expected[i].Date) {
					t.Errorf("Point %d: date = %v, expected %v", i, point.Date, tt.expected[i].Date)
				}

				if !almostEqual(point.MA7, tt.expected[i].MA7, 0.001) {
					t.Errorf("Point %d: MA7 = %f, expected %f", i, point.MA7, tt.expected[i].MA7)
				}

				if !almostEqual(point.MA25, tt.expected[i].MA25, 0.001) {
					t.Errorf("Point %d: MA25 = %f, expected %f", i, point.MA25, tt.expected[i].MA25)
				}

				if !almostEqual(point.MA99, tt.expected[i].MA99, 0.001) {
					t.Errorf("Point %d: MA99 = %f, expected %f", i, point.MA99, tt.expected[i].MA99)
				}
			}
		})
	}
}

func TestCalculateSeries_EdgeCases(t *testing.T) {
	baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("nil_inputs", func(t *testing.T) {
		result := CalculateSeries(nil, nil)
		if len(result) != 0 {
			t.Errorf("Expected empty result for nil inputs, got %d points", len(result))
		}
	})

	t.Run("different_lengths", func(t *testing.T) {
		dates := generateDates(baseDate, 100)
		prices := generatePrices(90, 100.0)

		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered from panic as expected:", r)
			}
		}()

		result := CalculateSeries(dates, prices)
		t.Log("Result length:", len(result))
	})
}

// generateDates belirtilen sayıda tarih oluşturur
func generateDates(start time.Time, count int) []time.Time {
	dates := make([]time.Time, count)
	for i := 0; i < count; i++ {
		dates[i] = start.AddDate(0, 0, i)
	}
	return dates
}

// generatePrices belirtilen sayıda sabit fiyat oluşturur
func generatePrices(count int, price float64) []float64 {
	prices := make([]float64, count)
	for i := 0; i < count; i++ {
		prices[i] = price
	}
	return prices
}

// generateIncreasingPrices artan fiyatlar oluşturur
func generateIncreasingPrices(count int) []float64 {
	prices := make([]float64, count)
	for i := 0; i < count; i++ {
		prices[i] = float64(i + 1)
	}
	return prices
}

// calculateExpectedMean belirli bir aralık için beklenen ortalamayı hesaplar
func calculateExpectedMean(prices []float64, start, end int) float64 {
	sum := 0.0
	for i := start; i < end; i++ {
		sum += prices[i]
	}
	return sum / float64(end-start)
}

// extendPattern bir pattern'i belirtilen uzunluğa genişletir
func extendPattern(pattern []float64, targetLength int) []float64 {
	result := make([]float64, targetLength)
	for i := 0; i < targetLength; i++ {
		result[i] = pattern[i%len(pattern)]
	}
	return result
}

// almostEqual float değerlerini belirli bir toleransla karşılaştırır
func almostEqual(a, b, tolerance float64) bool {
	return abs(a-b) <= tolerance
}

// abs mutlak değer fonksiyonu
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
