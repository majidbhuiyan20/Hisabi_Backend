package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string  `json:"name"`
	Barcode       string  `json:"barcode"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	Stock         int     `json:"stock"`
	Origin        string  `json:"origin"`
}
