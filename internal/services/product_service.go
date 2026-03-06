package services

import (
	"hisabi.com/m/internal/model"
	"hisabi.com/m/internal/repository"
)

func AddProduct(product *model.Product) error {
	return repository.CreateProduct(product)
}

func ListProduct() ([]model.Product, error) {
	return repository.GetAllProducts()
}

func UpdateProductService(id uint, updatedData *model.Product) (*model.Product, error) {
	return repository.UpdateProduct(id, updatedData)
}

func DeleteProductService(id uint) error {
	return repository.DeleteProduct(id)
}
