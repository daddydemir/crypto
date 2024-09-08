package config

import (
	"github.com/joho/godotenv"
	"os"
)

var path string

func init() {
	path = os.Getenv("ENV_PATH")
	if path == "" {
		path = "../config/envs/prod.env"
	}

	println("path: ", path)
}

func Get(key string) string {

	err := godotenv.Load(path)
	if err != nil {
		println("Error loading .env file", err)
	}
	return os.Getenv(key)
}
