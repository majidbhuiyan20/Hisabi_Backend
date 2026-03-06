package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

// Update Product Handler Code

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  false,
			Message: "Invalid product ID",
			Data:    nil,
		})
		return

	}

	var updatedProduct model.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	product, err := services.UpdateProductService(uint(id), &updatedProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// Delete Product by Id

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  false,
			Message: "Invalid product ID",
			Data:    nil,
		})
		return
	}

	// Call Service to delete product

	err = services.DeleteProductService(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  false,
			Message: "Failed to delete product: " + err.Error(),
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  true,
		Message: "Product Deleted Successfully",
		Data:    nil,
	})
}
