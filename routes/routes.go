package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"hisabi.com/m/internal/handler"
	"hisabi.com/m/middleware"
)

func SetUpRoutes() *mux.Router {
	router := mux.NewRouter()

	// ── Health Check ──────────────────────────────────────
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	api := router.PathPrefix("/api/v1").Subrouter()

	// ── Public Routes (no token needed) ──────────────────
	api.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	api.HandleFunc("/refresh", handler.RefreshHandler).Methods("POST")

	// Login has rate limit — 10 attempts per 5 minutes
	api.Handle("/login",
		middleware.LoginRateLimit(
			http.HandlerFunc(handler.LoginHandler),
		),
	).Methods("POST")

	// ── Protected Routes (JWT token required) ─────────────
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthRequired)

	protected.HandleFunc("/me", handler.MeHandler).Methods("GET")

	protected.HandleFunc("/products", handler.ProductHandler).Methods("GET", "POST")
	protected.HandleFunc("/products/{id}", handler.UpdateProductHandler).Methods("PUT")
	protected.HandleFunc("/products/{id}", handler.DeleteProductHandler).Methods("DELETE")

	return router
}

// func SetUpRoutes() *mux.Router {
// 	router := mux.NewRouter()
// 	api := router.PathPrefix("/api/v1").Subrouter()

// 	// ── Public ────────────────────────────────────────────
// 	api.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
// 	api.Handle("/login",
// 		middleware.LoginRateLimit(http.HandlerFunc(handler.LoginHandler)),
// 	).Methods("POST")
// 	api.HandleFunc("/refresh", handler.RefreshHandler).Methods("POST")

// 	api.HandleFunc("/products", handler.ProductHandler).Methods("GET", "POST")
// 	api.HandleFunc("/products/{id}", handler.UpdateProductHandler).Methods("PUT")
// 	api.HandleFunc("/products/{id}", handler.DeleteProductHandler).Methods("DELETE")

// 	// ── Protected ─────────────────────────────────────────
// 	protected := api.PathPrefix("").Subrouter()
// 	protected.Use(middleware.AuthRequired)
// 	protected.HandleFunc("/me", handler.MeHandler).Methods("GET")

// 	return router
// }
