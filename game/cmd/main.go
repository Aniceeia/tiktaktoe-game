package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game/internal/domain/services"
	"game/internal/infrastructure/repositories"
	httphandler "game/internal/interfaces/http"

	_ "github.com/lib/pq"
)

func main() {
	// Получаем переменные окружения
	dbURL := getEnv("DATABASE_URL", "postgres://game:password@localhost:5432/game_db?sslmode=disable")
	port := getEnv("PORT", "8080")

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение с БД
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализируем репозитории
	userRepo := repositories.NewUserRepository(db)
	gameRepo := repositories.NewGameRepository(db)

	// Инициализируем сервисы
	authService := services.NewAuthService(userRepo)
	aiService := services.NewAI()
	gameService := services.NewGameService(gameRepo, userRepo, aiService)

	// Инициализируем роутер
	router := httphandler.NewRouter(gameService, authService)

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router.GetRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
