package handler

import (
	"encoding/json"
	"net/http"

	"hisabi.com/m/internal/model"
	"hisabi.com/m/internal/services"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		products, err := services.ListProduct()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)

	case http.MethodPost:
		var product model.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := services.AddProduct(&product); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}
