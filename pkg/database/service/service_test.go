package service

import (
	"fmt"
	"github.com/daddydemir/crypto/config/database"
	"testing"
)

func TestCreateDaily(t *testing.T) {

	database.InitMySQLConnect()
	CreateDaily(false)
}

func TestGetDailyFromDatabase(t *testing.T) {
	database.InitMySQLConnect()
	models := GetDailyFromDatabase()
	//SortSlice(models)

	for i := 0; i < len(models); i++ {
		fmt.Println(models[i].ExchangeId)
	}
}
