package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, status bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(res)
}

func JSONStatus(w http.ResponseWriter, httpStatus int, status bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(res)
}
