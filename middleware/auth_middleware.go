package middleware

import (
	"context"
	"net/http"
	"strings"

	"hisabi.com/m/internal/services"
	"hisabi.com/m/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		authHeader := r.Header.Get("Authorization")

		if authHeader == ""{
			utils.JSONStatus(w, http.StatusUnauthorized, false, 
			"Authorization Header is missing", nil)
			return 
		}

		if !strings.HasPrefix(authHeader, "Bearer "){
			utils.JSONStatus(w, http.StatusUnauthorized, false, "Invalid token fromet. Use: Bearer <token>", nil)
			return 
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := services.VerifyAccessToken(tokenStr)

		if err != nil{
			utils.JSONStatus(w, http.StatusUnauthorized, false, "Token is invalid or has expried. Please login again", nil)
			return 
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request)uint{
	userID, _ := r.Context().Value(UserIDKey).(uint)
	return userID
}