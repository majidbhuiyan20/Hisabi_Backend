package handler

import (
	"encoding/json"
	"net/http"

	"hisabi.com/m/internal/model"
	"hisabi.com/m/internal/services"
	"hisabi.com/m/utils"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		products, err := services.ListProduct()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(utils.Response{
				Status:  false,
				Message: "Failed to fetch products: " + err.Error(),
				Data:    nil,
			})
			return
		}
		//json.NewEncoder(w).Encode(products)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  true,
			Message: "All Products Get Successfully",
			Data:    products,
		})

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
		// json.NewEncoder(w).Encode(product)

		json.NewEncoder(w).Encode(utils.Response{
			Status:  true,
			Message: "Product created successfully",
			Data:    product,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}
