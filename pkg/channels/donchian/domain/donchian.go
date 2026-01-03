package domain

import (
	"fmt"
	"time"
)

type DonchianChannel struct {
	Upper  float64
	Lower  float64
	Middle float64
	Date   string
	Price  float64
}

type DonchianData struct {
	Date  time.Time
	Min   float64
	Max   float64
	Close float64
}

func CalculateDonchian(data []DonchianData) ([]DonchianChannel, error) {

	if len(data) < 20 {
		return nil, fmt.Errorf("not enough data to calculate donchian channel, data length: %d", len(data))
	}

	result := make([]DonchianChannel, len(data)-20)
	for i := 0; i < len(data)-20; i++ {
		lowest, highest := getMinAndMaxWithPeriod(data[i:i+20], 20)
		result[i] = DonchianChannel{
			Upper:  highest,
			Lower:  lowest,
			Middle: (highest + lowest) / 2,
			Date:   data[i+20].Date.Format("2006-01-02"),
			Price:  data[i+20].Close,
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
