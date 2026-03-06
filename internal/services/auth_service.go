package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"hisabi.com/m/internal/model"
	"hisabi.com/m/internal/repository"
)

var jwtSecret = []byte("superSecretkey")

// hash Password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// Check Password
func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// REgister new user
func Register(username, email, password string) (*model.User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username: username,
		Email:    email,
		Password: hashed,
	}
	err = repository.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// generate jwt token

func GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtSecret)
}

// Login user

func Login(email, password string) (string, error) {
	user, err := repository.GetUserByemail(email)
	if err != nil {
		return "", errors.New("Invalid credentails")
	}
	if !checkPassword(user.Password, password) {
		return "", errors.New("Invalid Credentails")
	}
	return GenerateToken(user)
}
