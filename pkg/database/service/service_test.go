package service

import (
	"github.com/daddydemir/crypto/config/database"
	"testing"
)

func TestCreateDaily(t *testing.T) {

	database.InitMySQLConnect()
	CreateDaily(true)
}
