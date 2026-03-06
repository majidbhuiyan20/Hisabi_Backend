package routes

import (
	"github.com/gorilla/mux"
	"hisabi.com/m/internal/handler"
)

func SetUpRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/products", handler.ProductHandler).Methods("GET", "POST")
	router.HandleFunc("/api/v1/products/{id}", handler.UpdateProductHandler).Methods("PUT")
	return router
}
