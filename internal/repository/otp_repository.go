package repository

import (
	"time"

	database "hisabi.com/m/databases"
	"hisabi.com/m/internal/model"
)

func SaveOTP(email, code string, expiresAt time.Time) error {
	database.DB.Where("email = ?", email).Delete(&model.OTP{})

	otp := &model.OTP{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}
	return database.DB.Create(otp).Error
}
