package service

import (
	"fmt"
	"testing"
)

func TestGetWeek(t *testing.T) {
	//GetWeek()
	week, s := getWeek()
	fmt.Println(week, s)
}

func TestGetToday(t *testing.T) {
	today, s := getToday()
	fmt.Println(today, s)
}

func TestGetPeriodForTwoWeeks(t *testing.T) {
	start, end := getPeriodForTwoWeeks()
	fmt.Printf("start : %v, end : %v", start, end)

}
