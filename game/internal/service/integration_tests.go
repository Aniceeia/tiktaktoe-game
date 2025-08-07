package service

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://game:password@localhost/game_db_test?sslmode=disable")
	assert.NoError(t, err)

	// Создаем тестовые таблицы
	_, err = db.Exec(`
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS games;
		
		CREATE TABLE users (
			uuid VARCHAR(36) PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			score INT DEFAULT 0
		);
		
		CREATE TABLE games (
			id VARCHAR(36) PRIMARY KEY,
			board JSONB NOT NULL,
			status VARCHAR(50) NOT NULL,
			player1 VARCHAR(36) REFERENCES users(uuid) ON DELETE SET NULL,
			player2 VARCHAR(36),
			next_turn VARCHAR(36),
			mode VARCHAR(10) NOT NULL
		);
	`)
	assert.NoError(t, err)

	return db
}

func TestAuthService(t *testing.T) {
	db := setupDB(t)
	authService := NewPostgresAuthService(db)

	t.Run("Register and authenticate user", func(t *testing.T) {
		// Регистрация
		err := authService.Register(SignUpRequest{
			Login:    "testuser",
			Password: "password123",
		})
		assert.NoError(t, err)

		// Аутентификация
		uuid, err := authService.Authenticate("testuser", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, uuid)

		// Проверка получения пользователя
		user, err := authService.GetUserByUUID(uuid)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
	})
}

func TestGameService(t *testing.T) {
	db := setupDB(t)
	repo := NewPostgresGameStorage(db)
	ai := NewAI()
	gameService := NewGameService(repo, ai)
	authService := NewPostgresAuthService(db)

	// Создаем тестовых пользователей
	err := authService.Register(SignUpRequest{Login: "player1", Password: "p1"})
	assert.NoError(t, err)
	err = authService.Register(SignUpRequest{Login: "player2", Password: "p2"})
	assert.NoError(t, err)

	player1, _ := authService.Authenticate("player1", "p1")
	player2, _ := authService.Authenticate("player2", "p2")

	t.Run("Create and join game", func(t *testing.T) {
		// Создаем игру
		game, err := gameService.CreateGame(player1, "pvp")
		assert.NoError(t, err)
		assert.Equal(t, player1, game.Player1)
		assert.Equal(t, "", game.Player2)

		// Присоединяемся к игре
		err = gameService.JoinGame(game.ID, player2)
		assert.NoError(t, err)

		// Проверяем состояние игры
		updatedGame, err := gameService.GetGameState(game.ID)
		assert.NoError(t, err)
		assert.Equal(t, player2, updatedGame.Player2)
	})
}
