package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game/internal/domain/services"
	"game/internal/infrastructure/repositories"
	httphandler "game/internal/interfaces/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://game:password@localhost:5432/game_db?sslmode=disable")
	port := getEnv("PORT", "8080")

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	gameRepo := repositories.NewGameRepository(db)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userService)
	aiService := services.NewAI()
	gameService := services.NewGameService(gameRepo, userRepo, aiService)

	router := httphandler.NewRouter(gameService, authService)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router.GetRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
