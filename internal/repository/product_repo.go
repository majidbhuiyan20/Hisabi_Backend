package repository

import (
	database "hisabi.com/m/databases"
	"hisabi.com/m/internal/model"
)

func CreateProduct(product *model.Product) error {
	return database.DB.Create(product).Error
}
func GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	err := database.DB.Find(&products).Error
	return products, err
}
