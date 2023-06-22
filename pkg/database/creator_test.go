package database

import (
	"github.com/daddydemir/crypto/config/database"
	"testing"
)

func TestCreateTables(t *testing.T) {
	database.InitMySQLConnect()
	CreateTables()
}
