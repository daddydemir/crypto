package service

import (
	"github.com/daddydemir/crypto/pkg/model"
	"sort"
	"time"
)

func getWeek() (string, string) {
	today := time.Now()

	weekStart := today.AddDate(0, 0, int(time.Monday-today.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 6)

	return time2string(weekStart), time2string(weekEnd)
}

func getToday() (string, string) {
	today := time.Now()
	nextDay := today.AddDate(0, 0, 1)
	return today.Format("2006-01-02"), nextDay.Format("2006-01-02")
}

func sortSlice(arr []model.DailyModel) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].ExchangeId < arr[j].ExchangeId
	})
}

func sortSliceWithRate(arr []model.DailyModel) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Rate > arr[j].Rate
	})
}

func getPeriodForTwoWeeks() (string, string) {
	today := time.Now
	start := today().AddDate(0, 0, -14)

	return time2string(start), time2string(today())
}

func time2string(t time.Time) string {
	return t.Format("2006-01-02")
}
