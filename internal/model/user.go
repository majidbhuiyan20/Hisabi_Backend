package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex;not null;size:100"`
	Email    string `json:"email"    gorm:"uniqueIndex;not null;size:100"`
	Password string `json:"-"        gorm:"not null"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}
