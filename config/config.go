package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseURL string

	JWTAccessSecret  string
	JWTRefreshSecret string

	// SMTP fields মুছে দাও, Resend দিয়ে replace
	ResendAPIKey string // ← নতুন

	Port string
}

var Config *AppConfig

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, reading from environment")
	}

	Config = &AppConfig{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5433"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "1234"),
		DBName:      getEnv("DB_NAME", "hisabi"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		JWTAccessSecret:  mustGetEnv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: mustGetEnv("JWT_REFRESH_SECRET"),

		ResendAPIKey: mustGetEnv("RESEND_API_KEY"),

		Port: getEnv("PORT", "8080"),
	}

	log.Println("✅ Config loaded")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("❌ Required env variable missing: %s", key)
	}
	return v
}

func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}
