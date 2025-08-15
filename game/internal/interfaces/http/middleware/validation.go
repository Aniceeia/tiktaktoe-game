package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				respondWithError(w, "Invalid content type", "INVALID_CONTENT_TYPE", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		respondWithError(w, "Invalid JSON", "INVALID_JSON", http.StatusBadRequest)
		return false
	}

	if err := validate.Struct(v); err != nil {
		respondWithError(w, "Validation failed", "VALIDATION_FAILED", http.StatusBadRequest)
		return false
	}

	return true
}
