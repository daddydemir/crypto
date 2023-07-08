package service

import (
	"fmt"
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/pkg/dao"
	"github.com/daddydemir/crypto/pkg/model"
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

func TestGetDailyWithId(t *testing.T) {
	database.InitMySQLConnect()
	var date dao.Date
	date.Id = "gala"
	date.StartDate = "2023-06-22"
	date.EndDate = "2023-06-30"

	response := GetDailyWithId(date)
	fmt.Println(response)
}

func TestCreateMessage(t *testing.T) {
	database.InitMySQLConnect()
	CreateMessage()
}

func TestTemp(t *testing.T) {
	database.InitMySQLConnect()
	var date dao.Date
	date.StartDate, date.EndDate = getToday()
	db := GetDailyFromDb(date)
	m := make(map[string]model.DailyModel)

	for i := 0; i < len(db); i++ {
		m = MergeMap(m, db[i].ToMap())
	}
	fmt.Println(len(m))
}
