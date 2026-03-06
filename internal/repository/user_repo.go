package repository

import (
	database "hisabi.com/m/databases"
	"hisabi.com/m/internal/model"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func GetUserByemail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
