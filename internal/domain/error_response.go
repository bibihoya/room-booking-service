package domain

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Err Error `json:"error"`
}

func SendError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Err: Error{
			Code:    code,
			Message: message,
		},
	})
}

func SendSuccess(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
