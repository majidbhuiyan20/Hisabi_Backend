package handler

import (
	"encoding/json"
	"net/http"

	"hisabi.com/m/internal/services"
	"hisabi.com/m/middleware"
	"hisabi.com/m/utils"
)

type AuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

	utils.JSONStatus(w, http.StatusCreated, true, "Account created successfully",
		map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	)
}

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
		"access_token":  tokens.AccessToken,  // valid for 1 hour
		"refresh_token": tokens.RefreshToken, // valid for 30 days
	})
}

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

// GET /api/v1/me  — protected
func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	utils.JSON(w, true, "success", map[string]interface{}{
		"user_id": userID,
	})
}
