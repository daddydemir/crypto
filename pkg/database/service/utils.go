package service

import (
	"fmt"
	"time"
)

func getWeek() (string, string) {
	// Bugünün tarihini almak
	today := time.Now()
	fmt.Println("Bugün:", today.Format("2006-01-02"))

	// Bu haftanın başlangıcı ve bitişini belirlemek
	weekStart := today.AddDate(0, 0, int(time.Monday-today.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 6)

	//fmt.Println("Bu haftanın başlangıcı:", weekStart.Format("2006-01-02"))
	//fmt.Println("Bu haftanın bitişi:", weekEnd.Format("2006-01-02"))

	return weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")
}
