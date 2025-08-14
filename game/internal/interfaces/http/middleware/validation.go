package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationMiddleware проверяет валидность входящих запросов
func ValidationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем GET запросы
			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// Проверяем Content-Type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				respondWithError(w, errInvalidContentType, http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateRequest валидирует структуру запроса
func ValidateRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		respondWithError(w, errInvalidJSON, http.StatusBadRequest)
		return false
	}

	if err := validate.Struct(v); err != nil {
		respondWithError(w, errValidationFailed, http.StatusBadRequest)
		return false
	}

	return true
}
