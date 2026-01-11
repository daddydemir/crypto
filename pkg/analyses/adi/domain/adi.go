package domain

import "time"

type Point struct {
	Date time.Time `json:"date"`
	ADI  float64   `json:"adi"`
	PDI  float64   `json:"pdi"` // Positive Directional Indicator
	MDI  float64   `json:"mdi"` // Minus Directional Indicator
	DX   float64   `json:"dx"`  // Directional Movement Index
}

type PriceData struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// CalculateADI calculates Average Directional Index
func CalculateADI(prices []PriceData, period int) []Point {
	if len(prices) < period+1 {
		return nil
	}

	points := make([]Point, 0)

	// Calculate True Range, +DM, -DM for each day
	trueRanges := make([]float64, 0)
	plusDMs := make([]float64, 0)
	minusDMs := make([]float64, 0)

	for i := 1; i < len(prices); i++ {
		prev := prices[i-1]
		curr := prices[i]

		// True Range calculation
		tr1 := curr.High - curr.Low
		tr2 := abs(curr.High - prev.Close)
		tr3 := abs(curr.Low - prev.Close)
		tr := max(tr1, max(tr2, tr3))
		trueRanges = append(trueRanges, tr)

		// Directional Movement calculation
		upMove := curr.High - prev.High
		downMove := prev.Low - curr.Low

		var plusDM, minusDM float64
		if upMove > downMove && upMove > 0 {
			plusDM = upMove
		}
		if downMove > upMove && downMove > 0 {
			minusDM = downMove
		}

		plusDMs = append(plusDMs, plusDM)
		minusDMs = append(minusDMs, minusDM)
	}

	// Calculate smoothed averages
	if len(trueRanges) < period {
		return points
	}

	for i := period - 1; i < len(trueRanges); i++ {
		// Calculate ATR (smoothed TR)
		atr := calculateSmoothedAverage(trueRanges, i, period)

		// Calculate smoothed +DM and -DM
		smoothedPlusDM := calculateSmoothedAverage(plusDMs, i, period)
		smoothedMinusDM := calculateSmoothedAverage(minusDMs, i, period)

		// Calculate +DI and -DI
		pdi := (smoothedPlusDM / atr) * 100
		mdi := (smoothedMinusDM / atr) * 100

		// Calculate DX
		var dx float64
		if pdi+mdi != 0 {
			dx = (abs(pdi-mdi) / (pdi + mdi)) * 100
		}

		points = append(points, Point{
			Date: prices[i+1].Date,
			PDI:  pdi,
			MDI:  mdi,
			DX:   dx,
		})
	}

	// Calculate ADX (smoothed DX)
	if len(points) < period {
		return points
	}

	dxValues := make([]float64, len(points))
	for i, point := range points {
		dxValues[i] = point.DX
	}

	for i := period - 1; i < len(points); i++ {
		adx := calculateSmoothedAverage(dxValues, i, period)
		points[i].ADI = adx
	}

	// Return only points with calculated ADX
	return points[period-1:]
}

func calculateSmoothedAverage(values []float64, endIndex, period int) float64 {
	if endIndex < period-1 {
		return 0
	}

	if endIndex == period-1 {
		// First smoothed value is simple average
		sum := 0.0
		for i := endIndex - period + 1; i <= endIndex; i++ {
			sum += values[i]
		}
		return sum / float64(period)
	}

	// Subsequent values use Wilder's smoothing
	prevSmoothed := calculateSmoothedAverage(values, endIndex-1, period)
	return (prevSmoothed*(float64(period)-1) + values[endIndex]) / float64(period)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
