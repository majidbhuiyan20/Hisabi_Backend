package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string `json:"username"    gorm:"uniqueIndex;not null;size:100"`
	Email      string `json:"email"       gorm:"uniqueIndex;not null;size:100"`
	Password   string `json:"-"           gorm:"not null"`
	IsActive   bool   `json:"is_active"   gorm:"default:true"`
	IsVerified bool   `json:"is_verified" gorm:"default:false"`
}

type OTP struct {
	gorm.Model
	Email     string    `gorm:"not null;index;size:100"`
	Code      string    `gorm:"not null;size:6"`
	ExpiresAt time.Time `gorm:"not null"`
	IsUsed    bool      `gorm:"default:false"`
}
