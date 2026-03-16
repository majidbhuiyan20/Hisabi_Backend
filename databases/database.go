package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hisabi.com/m/config"
)

var DB *gorm.DB

func Connect() {
	var dsn string

	if config.Config.DatabaseURL != "" {
		dsn = config.Config.DatabaseURL
	} else {

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Config.DBHost,
			config.Config.DBPort,
			config.Config.DBUser,
			config.Config.DBPassword,
			config.Config.DBName,
		)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Database connected successfully")
}
