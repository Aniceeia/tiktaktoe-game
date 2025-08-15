package middleware

import (
	"context"
	"encoding/base64"
	"game/internal/domain/services"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldSkipAuth(r) {
				next.ServeHTTP(w, r)
				return
			}

			login, password, err := extractLoginPassword(r)
			if err != nil {
				respondWithError(w, "Invalid authorization header", "UNAUTHORIZED", http.StatusUnauthorized)
				return
			}

			// Валидация логина и пароля
			if login == "" || password == "" {
				respondWithError(w, "Invalid credentials", "UNAUTHORIZED", http.StatusUnauthorized)
				return
			}

			user, err := authService.Login(r.Context(), login, password)
			if err != nil || user == nil {
				respondWithError(w, "Invalid credentials", "UNAUTHORIZED", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user.UUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

func shouldSkipAuth(r *http.Request) bool {
	// Аутентификация не нужна если ты уже вошел
	if strings.Contains(r.URL.Path, "/auth/register") || strings.Contains(r.URL.Path, "/auth/login") {
		return true
	}
	return false
}

func extractLoginPassword(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", errInvalidHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", errInvalidHeader
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errInvalidHeader
	}

	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return "", "", errInvalidHeader
	}

	return credentials[0], credentials[1], nil
}

func respondWithError(w http.ResponseWriter, message, code string, statusCode int) {
	SendErrorResponse(w, message, code, statusCode)
}
