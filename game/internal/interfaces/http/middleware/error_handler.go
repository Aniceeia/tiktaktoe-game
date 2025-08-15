package middleware

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// SendErrorResponse - ошибки отформатированные под стандарт
func SendErrorResponse(w http.ResponseWriter, message, code string, statusCode int) {
	response := ErrorResponse{
		Error: message,
		Code:  code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func SendErrorResponseWithDetails(w http.ResponseWriter, message, code, details string, statusCode int) {
	response := ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
