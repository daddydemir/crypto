package service

import (
	"github.com/daddydemir/crypto/pkg/model"
	"sort"
	"time"
)

func getWeek() (string, string) {
	// Bugünün tarihini almak
	today := time.Now()

	// Bu haftanın başlangıcı ve bitişini belirlemek
	weekStart := today.AddDate(0, 0, int(time.Monday-today.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 6)

	//fmt.Println("Bu haftanın başlangıcı:", weekStart.Format("2006-01-02"))
	//fmt.Println("Bu haftanın bitişi:", weekEnd.Format("2006-01-02"))

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
	start := today().AddDate(0, 0, -13)

	return time2string(start), time2string(today())
}

func time2string(t time.Time) string {
	return t.Format("2006-01-02")
}
