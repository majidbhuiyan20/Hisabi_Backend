package repository

import (
	database "hisabi.com/m/databases"
	"hisabi.com/m/internal/model"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ? AND is_active = true", email).First(&user).Error
	return &user, err
}

func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.Where("id = ? AND is_active = true", id).First(&user).Error
	return &user, err
}

func EmailExists(email string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func UsernameExists(username string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
