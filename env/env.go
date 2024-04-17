package env

import (
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() error {
	return godotenv.Load(".env")
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
