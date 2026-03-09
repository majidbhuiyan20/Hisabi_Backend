package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"hisabi.com/m/internal/model"
	"hisabi.com/m/internal/services"
	"hisabi.com/m/middleware"
	"hisabi.com/m/utils"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// GET  /api/v1/products
// POST /api/v1/products
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func ProductHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	switch r.Method {

	case http.MethodGet:
		products, err := services.ListProduct(userID)
		if err != nil {
			utils.JSONStatus(w, http.StatusInternalServerError,
				false, "Failed to fetch products: "+err.Error(), nil)
			return
		}
		utils.JSON(w, true, "Products fetched successfully", products)

	case http.MethodPost:
		var product model.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			utils.JSONStatus(w, http.StatusBadRequest,
				false, "Invalid request body", nil)
			return
		}

		if err := services.AddProduct(userID, &product); err != nil {
			utils.JSONStatus(w, http.StatusBadRequest,
				false, err.Error(), nil)
			return
		}

		utils.JSONStatus(w, http.StatusCreated,
			true, "Product created successfully", product)

	default:
		utils.JSONStatus(w, http.StatusMethodNotAllowed,
			false, "Method not allowed", nil)
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// PUT /api/v1/products/{id}
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONStatus(w, http.StatusBadRequest,
			false, "Invalid product ID", nil)
		return
	}

	var updatedProduct model.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		utils.JSONStatus(w, http.StatusBadRequest,
			false, "Invalid request body", nil)
		return
	}

	product, err := services.UpdateProductService(uint(id), userID, &updatedProduct)
	if err != nil {
		utils.JSONStatus(w, http.StatusNotFound,
			false, err.Error(), nil)
		return
	}

	utils.JSON(w, true, "Product updated successfully", product)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// DELETE /api/v1/products/{id}
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.JSONStatus(w, http.StatusBadRequest,
			false, "Invalid product ID", nil)
		return
	}

	if err := services.DeleteProductService(uint(id), userID); err != nil {
		utils.JSONStatus(w, http.StatusNotFound,
			false, err.Error(), nil)
		return
	}

	utils.JSON(w, true, "Product deleted successfully", nil)
}
