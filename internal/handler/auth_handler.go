package handler

import (
	"encoding/json"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, false, "Invalid input", nil)
		return
	}
}
