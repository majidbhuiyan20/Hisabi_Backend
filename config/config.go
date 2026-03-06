package config

import (
	"os"
)

var (
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
)

func Load() {
	DBHost = getEnv("DB_HOST", "loalhost")
	DBPort = getEnv("DB_PORT", "5433")
	DBUser = getEnv("DB_USER", "postgres")
	DBPassword = getEnv("DP_PASSWORD", "1234")
	DBName = getEnv("DB_NAME", "hisabi")
}
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
