package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"game/internal/domain/entities"
	"game/internal/domain/services"
	infrarepos "game/internal/infrastructure/repositories"
	"game/internal/interfaces/http/dto"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeaderboardIntegration(t *testing.T) {
	// Setup database connection
	ctx := context.Background()
	db, err := pgxpool.New(ctx, "postgres://game:game@localhost:5432/game_db?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	// Clean database
	_, err = db.Exec(ctx, "TRUNCATE users, games CASCADE")
	require.NoError(t, err)

	// Setup repositories and services
	userRepo := infrarepos.NewUserRepository(db)
	gameRepo := infrarepos.NewGameRepository(db)
	aiService := services.NewAIService()
	gameService := services.NewGameService(gameRepo, userRepo, aiService)

	// Create test users
	user1, err := userRepo.Create(ctx, "player1", "password123")
	require.NoError(t, err)

	user2, err := userRepo.Create(ctx, "player2", "password123")
	require.NoError(t, err)

	user3, err := userRepo.Create(ctx, "player3", "password123")
	require.NoError(t, err)

	// Create and play games to generate stats
	t.Run("Create games with different outcomes", func(t *testing.T) {
		// Game 1: player1 wins
		game1 := entities.NewGame(user1.UUID, entities.ModePvP)
		err := gameRepo.Create(ctx, game1)
		require.NoError(t, err)

		// Join game
		err = game1.JoinGame(user2.UUID)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game1)
		require.NoError(t, err)

		// Play moves - player1 wins horizontally
		err = game1.MakeMove(user1.UUID, 0, 0)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game1)
		require.NoError(t, err)

		err = game1.MakeMove(user2.UUID, 1, 1)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game1)
		require.NoError(t, err)

		err = game1.MakeMove(user1.UUID, 0, 1)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game1)
		require.NoError(t, err)

		err = game1.MakeMove(user2.UUID, 2, 2)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game1)
		require.NoError(t, err)

		err = game1.MakeMove(user1.UUID, 0, 2)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game1)
		require.NoError(t, err)

		assert.Equal(t, entities.StatusPlayer1Won, game1.Status)

		// Game 2: player2 wins
		game2 := entities.NewGame(user2.UUID, entities.ModePvP)
		err = gameRepo.Create(ctx, game2)
		require.NoError(t, err)

		err = game2.JoinGame(user3.UUID)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game2)
		require.NoError(t, err)

		// Play moves - player2 wins vertically
		err = game2.MakeMove(user2.UUID, 0, 0)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game2)
		require.NoError(t, err)

		err = game2.MakeMove(user3.UUID, 1, 1)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game2)
		require.NoError(t, err)

		err = game2.MakeMove(user2.UUID, 1, 0)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game2)
		require.NoError(t, err)

		err = game2.MakeMove(user3.UUID, 2, 2)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game2)
		require.NoError(t, err)

		err = game2.MakeMove(user2.UUID, 2, 0)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game2)
		require.NoError(t, err)

		assert.Equal(t, entities.StatusPlayer2Won, game2.Status)

		// Game 3: player1 wins again
		game3 := entities.NewGame(user1.UUID, entities.ModePvP)
		err = gameRepo.Create(ctx, game3)
		require.NoError(t, err)

		err = game3.JoinGame(user3.UUID)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game3)
		require.NoError(t, err)

		// Play moves - player1 wins diagonally
		err = game3.MakeMove(user1.UUID, 0, 0)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game3)
		require.NoError(t, err)

		err = game3.MakeMove(user3.UUID, 0, 1)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game3)
		require.NoError(t, err)

		err = game3.MakeMove(user1.UUID, 1, 1)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game3)
		require.NoError(t, err)

		err = game3.MakeMove(user3.UUID, 0, 2)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game3)
		require.NoError(t, err)

		err = game3.MakeMove(user1.UUID, 2, 2)
		require.NoError(t, err)
		err = gameRepo.Update(ctx, game3)
		require.NoError(t, err)

		assert.Equal(t, entities.StatusPlayer1Won, game3.Status)
	})

	t.Run("Test leaderboard retrieval", func(t *testing.T) {
		// Get leaderboard
		players, err := gameService.GetLeaderboard(ctx, 10)
		require.NoError(t, err)
		require.Len(t, players, 3)

		// Find players by login
		var player1, player2, player3 *entities.LeaderboardPlayer
		for _, p := range players {
			switch p.Login {
			case "player1":
				player1 = p
			case "player2":
				player2 = p
			case "player3":
				player3 = p
			}
		}

		require.NotNil(t, player1)
		require.NotNil(t, player2)
		require.NotNil(t, player3)

		// Verify stats
		assert.Equal(t, 2, player1.Wins)
		assert.Equal(t, 0, player1.Losses)
		assert.Equal(t, 0, player1.Draws)
		assert.Equal(t, 2, player1.TotalGames)
		assert.Equal(t, 1.0, player1.WinRatio)

		assert.Equal(t, 1, player2.Wins)
		assert.Equal(t, 1, player2.Losses)
		assert.Equal(t, 0, player2.Draws)
		assert.Equal(t, 2, player2.TotalGames)
		assert.Equal(t, 0.5, player2.WinRatio)

		assert.Equal(t, 0, player3.Wins)
		assert.Equal(t, 2, player3.Losses)
		assert.Equal(t, 0, player3.Draws)
		assert.Equal(t, 2, player3.TotalGames)
		assert.Equal(t, 0.0, player3.WinRatio)

		// Verify ordering (by win ratio descending)
		assert.Equal(t, "player1", players[0].Login)
		assert.Equal(t, "player2", players[1].Login)
		assert.Equal(t, "player3", players[2].Login)
	})

	t.Run("Test leaderboard with limit", func(t *testing.T) {
		// Get top 2 players
		players, err := gameService.GetLeaderboard(ctx, 2)
		require.NoError(t, err)
		require.Len(t, players, 2)

		// Should return top 2 players
		assert.Equal(t, "player1", players[0].Login)
		assert.Equal(t, "player2", players[1].Login)
	})
}

func TestLeaderboardAPI(t *testing.T) {
	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:8080"

	t.Run("Test leaderboard API endpoint", func(t *testing.T) {
		// Register and login user
		registerResp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json",
			strings.NewReader(`{"login":"testuser","password":"password123"}`))
		require.NoError(t, err)
		defer registerResp.Body.Close()

		loginResp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json",
			strings.NewReader(`{"login":"testuser","password":"password123"}`))
		require.NoError(t, err)
		defer loginResp.Body.Close()

		var loginData dto.JwtResponse
		err = json.NewDecoder(loginResp.Body).Decode(&loginData)
		require.NoError(t, err)

		// Test leaderboard endpoint
		req, err := http.NewRequest("GET", baseURL+"/api/v1/leaderboard?limit=5", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+loginData.AccessToken)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboard dto.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboard)
		require.NoError(t, err)

		// Should return valid leaderboard structure
		assert.NotNil(t, leaderboard.Players)
		assert.GreaterOrEqual(t, leaderboard.Total, 0)
	})

	t.Run("Test unauthorized access", func(t *testing.T) {
		// Test without authorization header
		resp, err := http.Get(baseURL + "/api/v1/leaderboard")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
