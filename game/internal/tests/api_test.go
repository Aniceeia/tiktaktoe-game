package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"game/internal/api"
	"game/internal/service"

	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T) *httptest.Server {
	mockGameService := &service.MockGameService{
		CreateGameFunc: func(creator string, mode string) (*service.Game, error) {
			return &service.Game{
				ID:       "test-game",
				Player1:  creator,
				Mode:     mode,
				Status:   service.StatusWaiting,
				Board:    [3][3]int{},
				NextTurn: creator,
			}, nil
		},
		JoinGameFunc: func(gameID string, joiner string) error {
			return nil
		},
		ProcessMoveFunc: func(gameID string, user string, newBoard [3][3]int) (*service.Game, error) {
			return &service.Game{
				ID:       gameID,
				Player1:  "alice-uuid",
				Player2:  "bob-uuid",
				Board:    newBoard,
				Status:   service.StatusPlayer2Turn,
				NextTurn: "bob-uuid",
				Mode:     "pvp",
			}, nil
		},
	}

	mockAuthService := &service.MockAuthService{
		RegisterFunc: func(req service.SignUpRequest) error {
			return nil
		},
		AuthenticateFunc: func(login, password string) (string, error) {
			if login == "alice" {
				return "alice-uuid", nil
			} else if login == "bob" {
				return "bob-uuid", nil
			}
			return "user-uuid", nil
		},
		GetUserByUUIDFunc: func(uuid string) (*service.User, error) {
			return &service.User{
				UUID:     uuid,
				Username: "user-" + uuid,
			}, nil
		},
	}

	mockAI := &service.AI{}

	handler := api.NewGameHandler(mockGameService, mockAI, mockAuthService)
	router := api.NewRouter(handler)

	return httptest.NewServer(router)
}

func TestGameFlow(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Регистрация игроков
	players := []struct {
		username string
		password string
	}{
		{"alice", "alicepass"},
		{"bob", "bobpass"},
	}

	for _, player := range players {
		regData := map[string]string{
			"login":    player.username,
			"password": player.password,
		}
		regBody, _ := json.Marshal(regData)
		_, err := http.Post(server.URL+"/register", "application/json", bytes.NewBuffer(regBody))
		assert.NoError(t, err)
	}

	// Alice создает игру
	createReq, _ := http.NewRequest("POST", server.URL+"/games", bytes.NewBufferString(`{"mode":"pvp"}`))
	createReq.SetBasicAuth("alice", "alicepass")
	createResp, err := http.DefaultClient.Do(createReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)

	var game api.GameResponse
	json.NewDecoder(createResp.Body).Decode(&game)
	gameID := game.ID

	// Bob присоединяется к игре
	joinReq, _ := http.NewRequest("POST", server.URL+"/games/"+gameID+"/join", nil)
	joinReq.SetBasicAuth("bob", "bobpass")
	joinResp, err := http.DefaultClient.Do(joinReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, joinResp.StatusCode)

	// Alice делает ход
	moveData := map[string][3][3]int{
		"board": {
			{1, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	}
	moveBody, _ := json.Marshal(moveData)
	moveReq, _ := http.NewRequest("POST", server.URL+"/games/"+gameID+"/move", bytes.NewBuffer(moveBody))
	moveReq.SetBasicAuth("alice", "alicepass")
	moveResp, err := http.DefaultClient.Do(moveReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, moveResp.StatusCode)
}
