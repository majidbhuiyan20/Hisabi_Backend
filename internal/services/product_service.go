package services

import (
	"errors"

	"hisabi.com/m/internal/model"
	"hisabi.com/m/internal/repository"
)

func AddProduct(userID uint, product *model.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.SellingPrice <= 0 {
		return errors.New("selling price must be greater than 0")
	}
	if product.PurchasePrice <= 0 {
		return errors.New("purchase price must be greater than 0")
	}
	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	product.UserID = userID

	return repository.CreateProduct(product)
}

func ListProduct(userID uint) ([]model.Product, error) {
	return repository.GetAllProducts(userID)
}

func UpdateProductService(id uint, userID uint, updatedData *model.Product) (*model.Product, error) {
	if updatedData.Name == "" {
		return nil, errors.New("product name is required")
	}
	if updatedData.SellingPrice <= 0 {
		return nil, errors.New("selling price must be greater than 0")
	}
	return repository.UpdateProduct(id, userID, updatedData)
}

func DeleteProductService(id uint, userID uint) error {
	return repository.DeleteProduct(id, userID)
}
