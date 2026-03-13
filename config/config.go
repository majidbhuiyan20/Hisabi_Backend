package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTAccessSecret  string
	JWTRefreshSecret string

	// ← নতুন যোগ হয়েছে
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string

	Port string
}

var Config *AppConfig

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, reading from environment")
	}

	Config = &AppConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "1234"),
		DBName:     getEnv("DB_NAME", "hisabi"),

		JWTAccessSecret:  mustGetEnv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: mustGetEnv("JWT_REFRESH_SECRET"),

		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPEmail:    mustGetEnv("SMTP_EMAIL"),
		SMTPPassword: mustGetEnv("SMTP_PASSWORD"),

		Port: getEnv("PORT", "8080"),
	}

	log.Println("Config loaded")
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
		log.Fatalf("Required env variable missing: %s", key)
	}
	return v
}

func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}
