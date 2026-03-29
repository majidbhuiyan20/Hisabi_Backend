package handler

import (
	"encoding/json"
	"net/http"

	"hisabi.com/m/internal/services"
	"hisabi.com/m/middleware"
	"hisabi.com/m/utils"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// POST /api/v1/register
//
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, false, "Invalid request body", nil)
		return
	}

	user, err := services.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.JSON(w, false, err.Error(), nil)
		return
	}

	// ✅ Account বানানো হয়েছে, OTP গেছে email এ
	utils.JSONStatus(w, http.StatusCreated, true,
		"Account created! Please check your email for the verification code.",
		map[string]interface{}{
			"user_id":  user.ID,
			"email":    user.Email,
			"verified": false,
		},
	)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// POST /api/v1/verify-otp
//
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, false, "Invalid request body", nil)
		return
	}

	if err := services.VerifyOTP(req.Email, req.OTP); err != nil {
		utils.JSON(w, false, err.Error(), nil)
		return
	}

	utils.JSON(w, true, "Email verified successfully! You can now log in.", nil)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// POST /api/v1/resend-otp
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func ResendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, false, "Invalid request body", nil)
		return
	}

	if req.Email == "" {
		utils.JSON(w, false, "Email is required", nil)
		return
	}

	if err := services.ResendOTP(req.Email); err != nil {
		utils.JSON(w, false, err.Error(), nil)
		return
	}

	utils.JSON(w, true, "A new verification code has been sent to your email.", nil)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// POST /api/v1/login
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, false, "Invalid request body", nil)
		return
	}

	tokens, err := services.Login(req.Email, req.Password)
	if err != nil {
		utils.JSONStatus(w, http.StatusUnauthorized, false, err.Error(), nil)
		return
	}

	utils.JSON(w, true, "Login successful", map[string]interface{}{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// POST /api/v1/refresh
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, false, "Invalid request body", nil)
		return
	}
	if req.RefreshToken == "" {
		utils.JSON(w, false, "refresh_token is required", nil)
		return
	}

	newToken, err := services.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		utils.JSONStatus(w, http.StatusUnauthorized, false, err.Error(), nil)
		return
	}

	utils.JSON(w, true, "Token refreshed successfully", map[string]string{
		"access_token": newToken,
	})
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// GET /api/v1/me  — protected
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	utils.JSON(w, true, "success", map[string]interface{}{
		"user_id": userID,
	})
}
