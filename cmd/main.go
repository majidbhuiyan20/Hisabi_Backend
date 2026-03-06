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
	database.Connect()

	database.DB.AutoMigrate(&model.Product{})

	mux := routes.SetUpRoutes()

	log.Println("Server running on port 8080....")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
