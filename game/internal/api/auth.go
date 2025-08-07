package api

import (
	"context"
	"net/http"
)

type contextKey string

const (
	userUUIDKey = contextKey("userUUID")
)

func (h *GameHandler) UserAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем публичные эндпоинты
		if r.URL.Path == "/register" || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Получаем Basic Auth
		username, password, ok := r.BasicAuth()
		if !ok {
			respondError(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		uuid, err := h.authService.Authenticate(username, password)
		if err != nil {
			respondError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userUUIDKey, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
