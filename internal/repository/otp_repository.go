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

// ------------------------------
// Otp Search and Verify
// ------------------------------
func GetValidOTP(email, code string) (*model.OTP, error) {
	var otp model.OTP
	err := database.DB.Where(
		"email = ? AND code = ? AND expires_at > ? AND is_used = false",
		email, code, time.Now(),
	).First(&otp).Error
	return &otp, err

}

//------------------------------
// OTP mark as used
//------------------------------

func MarkOTPUsed(id uint) error {
	return database.DB.Model(&model.OTP{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

// ------------------------------
// User verified mark
// ------------------------------
func MarkUserVerified(email string) error {
	return database.DB.Model(&model.User{}).Where("email = ?", email).Update("is_verified", true).Error
}
