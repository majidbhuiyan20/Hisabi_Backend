package main

import (
	"log"

	"hisabi.com/m/config"
	database "hisabi.com/m/databases"
)

func main() {

	config.Load()
	database.Connect()

	log.Println("Server running on port 8080....")
	//log.Fatal(http.ListenAndServe(":8080", mux))
}
