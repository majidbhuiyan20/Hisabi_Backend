package services

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"hisabi.com/m/config"
	"hisabi.com/m/internal/model"
	"hisabi.com/m/internal/repository"
	"hisabi.com/m/utils"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Password
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Token Generation
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func generateAccessToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"token_type": tokenTypeAccess,
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.JWTAccessSecret))
}

func generateRefreshToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"token_type": tokenTypeRefresh,
		"exp":        time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.JWTRefreshSecret))
}

func generateTokenPair(user *model.User) (*TokenPair, error) {
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Register
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func Register(username, email, password string) (*model.User, error) {

	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))

	if err := utils.ValidateRegister(username, email, password); err != nil {
		return nil, err
	}

	exists, err := repository.EmailExists(email)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	if exists {
		return nil, errors.New("an account with this email already exists")
	}

	usernameExists, err := repository.UsernameExists(username)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	if usernameExists {
		return nil, errors.New("this username is already taken")
	}

	hashed, err := hashPassword(password)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: hashed,
		IsActive: true,
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, errors.New("failed to create account")
	}

	return user, nil
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Login
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func Login(email, password string) (*TokenPair, error) {

	email = strings.TrimSpace(strings.ToLower(email))

	if err := utils.ValidateLogin(email, password); err != nil {
		return nil, err
	}

	user, err := repository.GetUserByEmail(email)
	if err != nil {
		// Same error for wrong email or wrong password
		// Prevents user enumeration attacks
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("this account has been deactivated")
	}

	if !checkPassword(user.Password, password) {
		return nil, errors.New("invalid email or password")
	}

	return generateTokenPair(user)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Refresh Access Token
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func RefreshAccessToken(refreshTokenStr string) (string, error) {

	token, err := jwt.Parse(refreshTokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token method")
		}
		return []byte(config.Config.JWTRefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}

	if claims["token_type"] != tokenTypeRefresh {
		return "", errors.New("this is not a refresh token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	user, err := repository.GetUserByID(uint(userIDFloat))
	if err != nil || !user.IsActive {
		return "", errors.New("user not found")
	}

	return generateAccessToken(user)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Verify Access Token — used in middleware
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func VerifyAccessToken(tokenStr string) (uint, error) {

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.Config.JWTAccessSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token")
	}

	if claims["token_type"] != tokenTypeAccess {
		return 0, errors.New("this is not an access token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id in token")
	}

	return uint(userIDFloat), nil
}
