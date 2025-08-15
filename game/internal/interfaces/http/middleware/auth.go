package middleware

import (
	"context"
	"net/http"
	"strings"

	"game/internal/domain/services"
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

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, "Invalid authorization header", "UNAUTHORIZED", http.StatusUnauthorized)
				return
			}
			token := parts[1]
			userID, err := authService.ValidateAccessToken(token)
			if err != nil || userID == "" {
				respondWithError(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
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
	if strings.Contains(r.URL.Path, "/auth/register") ||
		strings.Contains(r.URL.Path, "/auth/login") ||
		strings.Contains(r.URL.Path, "/auth/refresh/access") {
		return true
	}
	return false
}

func respondWithError(w http.ResponseWriter, message, code string, statusCode int) {
	SendErrorResponse(w, message, code, statusCode)
}
