package handler

import (
	"encoding/json"
	"net/http"

	"hisabi.com/m/internal/services"
	"hisabi.com/m/utils"
)

type AuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSON(w, false, "Method not allowed", nil)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, false, "Invalid Input", nil)
		return
	}
	user, err := services.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.JSON(w, false, err.Error(), nil)
		return
	}
	utils.JSON(w, true, "User registered Successfully", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})
}

func LogicHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.JSON(w, false, "Method not allowed", nil)
		return
	}
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, false, "Invalid Input", nil)
		return
	}
	token, err := services.Login(req.Email, req.Password)

	if err != nil {
		utils.JSON(w, false, err.Error(), nil)
		return
	}

	utils.JSON(w, true, "Login Successful", map[string]interface{}{
		"email":    req.Email,
		"token":    token,
	})
}
