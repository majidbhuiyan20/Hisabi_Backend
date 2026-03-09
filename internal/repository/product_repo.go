package repository

import (
	"errors"
	"fmt"

	database "hisabi.com/m/databases"
	"hisabi.com/m/internal/model"
)

func CreateProduct(product *model.Product) error {
	return database.DB.Create(product).Error
}
func GetAllProducts(userID uint) ([]model.Product, error) {
	var products []model.Product
	err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// Update Products by ID

func UpdateProduct(id uint, userID uint, updatedData *model.Product) (*model.Product, error) {
	var product model.Product

	err := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&product).Error
	if err != nil {
		return nil, errors.New("product not found or you do not have permission")
	}

	product.Name = updatedData.Name
	product.Barcode = updatedData.Barcode
	product.PurchasePrice = updatedData.PurchasePrice
	product.SellingPrice = updatedData.SellingPrice
	product.Stock = updatedData.Stock
	product.Origin = updatedData.Origin

	if err := database.DB.Save(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Delete Product by ID

func DeleteProduct(id uint, userID uint) error {
	result := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Product{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found or you do not have permission")
	}
	return nil
}
