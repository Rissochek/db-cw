package utils

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadEnvFile() {
	if err := godotenv.Load(); err != nil {
		zap.S().Fatalf("no .env file found")
	}
}

func GetKeyFromEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		zap.S().Fatalf("value is not set in .env file: %v", key)
	}
	return value
}