package config

import (
	"os"

	"github.com/joho/godotenv"
)

var path string

func init() {
	path = ".env"
}

func Get(key string) string {

	err := godotenv.Load(path)
	if err != nil {
		println("Error loading .env file", err)
	}
	return os.Getenv(key)
}
