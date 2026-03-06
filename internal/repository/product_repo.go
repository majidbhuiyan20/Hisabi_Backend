package repository

import (
	"errors"

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

// Update Products by ID

func UpdateProduct(id uint, updatedData *model.Product) (*model.Product, error) {
	var product model.Product

	err := database.DB.First(&product, id).Error
	if err != nil {
		return nil, errors.New("Product not found")
	}

	product.Name = updatedData.Name
	product.Barcode = updatedData.Barcode
	product.PurchasePrice = updatedData.PurchasePrice
	product.SellingPrice = updatedData.SellingPrice
	product.Stock = updatedData.Stock
	product.Origin = updatedData.Origin

	err = database.DB.Save(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
