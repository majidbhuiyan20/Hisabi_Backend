package main

import (
	"log"
	"net/http"

	"hisabi.com/m/config"
	database "hisabi.com/m/databases"
	"hisabi.com/m/internal/model"
	"hisabi.com/m/routes"
)

func main() {

	config.Load()

	// ② Database connect
	database.Connect()

	err := database.DB.AutoMigrate(
		&model.User{},
		&model.OTP{},
		&model.Product{},
	)
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	log.Println("Migration complete")

	mux := routes.SetUpRoutes()

	port := config.Config.Port
	log.Printf("🚀 Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
