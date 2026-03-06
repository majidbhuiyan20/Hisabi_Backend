package routes

import (
	"net/http"

	"hisabi.com/m/internal/handler"
)

func SetUpRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/products", handler.ProductHandler)
	return mux
}
