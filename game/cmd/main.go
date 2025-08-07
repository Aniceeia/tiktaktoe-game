package main

import (
	"database/sql"
	"game/internal/api"
	"game/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Проверка соединения с БД
	err = db.Ping()
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}
	log.Println("Successfully connected to database")

	repo := service.NewPostgresGameStorage(db)
	authService := service.NewPostgresAuthService(db)

	ai := service.NewAI()
	svc := service.NewGameService(repo, ai)

	handler := api.NewGameHandler(svc, ai, authService)
	router := api.NewRouter(handler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	server.Close()
}
