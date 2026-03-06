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
