package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware логирует HTTP запросы
func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Создаем wrapper для ResponseWriter, чтобы перехватить статус код
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			log.Printf(
				"%s %s %s %d %v",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				wrapped.statusCode,
				duration,
			)
		})
	}
}

// responseWriter обертка для http.ResponseWriter для перехвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
