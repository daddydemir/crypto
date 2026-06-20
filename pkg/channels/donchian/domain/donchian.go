package domain

import (
	"fmt"
	"time"
)

type DonchianChannel struct {
	Upper     float64
	Lower     float64
	Middle    float64
	Date      string
	Price     float64
	LowPrice  float64
	HighPrice float64
}

type DonchianData struct {
	Date  time.Time
	Min   float64
	Max   float64
	Close float64
}

func CalculateDonchian(data []DonchianData, period int) ([]DonchianChannel, error) {

	if len(data) < period {
		return nil, fmt.Errorf("not enough data to calculate donchian channel, data length: %d", len(data))
	}

	result := make([]DonchianChannel, len(data)-period)
	for i := 0; i < len(data)-period; i++ {
		lowest, highest := getMinAndMaxWithPeriod(data[i:i+period], period)
		result[i] = DonchianChannel{
			Upper:     highest,
			Lower:     lowest,
			Middle:    (highest + lowest) / 2,
			Date:      data[i+period].Date.Format("2006-01-02"),
			Price:     data[i+period].Close,
			LowPrice:  data[i+period].Min,
			HighPrice: data[i+period].Max,
		}
	}

	return result, nil
}

func getMinAndMaxWithPeriod(array []DonchianData, period int) (float64, float64) {
	lowest, highest := array[0].Min, array[0].Max
	for i := 1; i < period; i++ {
		if array[i].Min < lowest {
			lowest = array[i].Min
		}
		if array[i].Max > highest {
			highest = array[i].Max
		}
	}
	return lowest, highest
}
