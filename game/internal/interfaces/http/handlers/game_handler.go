package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"game/internal/domain/entities"
	"game/internal/domain/services"
	"game/internal/interfaces/http/dto"
	"game/internal/interfaces/http/middleware"

	"github.com/gorilla/mux"
)

type GameHandler struct {
	gameService *services.GameService
	authService *services.AuthService
}

func NewGameHandler(gameService *services.GameService, authService *services.AuthService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		authService: authService,
	}
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	mode := entities.GameMode(req.Mode)
	if mode != entities.ModePvP && mode != entities.ModePvE {
		h.sendErrorResponse(w, "Invalid game mode", "INVALID_MODE", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	game, err := h.gameService.CreateGame(r.Context(), userID, mode)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	var player2 *entities.User
	if game.Player2ID != "" && game.Player2ID != "AI" {
		player2, err = h.authService.GetUserByID(r.Context(), game.Player2ID)
		if err != nil {
			h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
	}

	response := dto.ToGameResponse(game, player1, player2)
	h.sendJSONResponse(w, response, http.StatusCreated)
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	err := h.gameService.JoinGame(r.Context(), gameID, userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	game, err := h.gameService.GetGameState(r.Context(), gameID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	var player2 *entities.User
	if game.Player2ID != "" && game.Player2ID != "AI" {
		player2, err = h.authService.GetUserByID(r.Context(), game.Player2ID)
		if err != nil {
			h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
	}

	response := dto.ToGameResponse(game, player1, player2)
	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	var req dto.MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	game, err := h.gameService.MakeMove(r.Context(), gameID, userID, req.Row, req.Col)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	var player2 *entities.User
	if game.Player2ID != "" && game.Player2ID != "AI" {
		player2, err = h.authService.GetUserByID(r.Context(), game.Player2ID)
		if err != nil {
			h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
	}

	response := dto.ToGameResponse(game, player1, player2)
	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *GameHandler) GetGameState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	game, err := h.gameService.GetGameState(r.Context(), gameID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	if !game.IsPlayerInGame(userID) {
		h.sendErrorResponse(w, "Forbidden", "FORBIDDEN", http.StatusForbidden)
		return
	}

	player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	var player2 *entities.User
	if game.Player2ID != "" && game.Player2ID != "AI" {
		player2, err = h.authService.GetUserByID(r.Context(), game.Player2ID)
		if err != nil {
			h.sendErrorResponse(w, "Failed to get player info", "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
	}

	response := dto.ToGameResponse(game, player1, player2)
	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *GameHandler) GetAvailableGames(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	games, err := h.gameService.GetAvailableGames(r.Context())
	if err != nil {
		handleError := err
		_ = handleError
		h.handleServiceError(w, err)
		return
	}

	var gameInfos []dto.GameShortInfo
	for _, game := range games {
		player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
		if err != nil {
			continue
		}

		var player2 *entities.User
		if game.Player2ID != "" && game.Player2ID != "AI" {
			player2, err = h.authService.GetUserByID(r.Context(), game.Player2ID)
			if err != nil {
				continue
			}
		}

		gameInfo := dto.ToGameShortInfo(game, player1, player2)
		gameInfos = append(gameInfos, gameInfo)
	}

	response := dto.AvailableGamesResponse{
		Games: gameInfos,
		Total: len(gameInfos),
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *GameHandler) GetUserGames(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	games, err := h.gameService.GetUserGames(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	var gameInfos []dto.GameShortInfo
	for _, game := range games {
		player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
		if err != nil {
			continue
		}

		var player2 *entities.User
		if game.Player2ID != "" && game.Player2ID != "AI" {
			player2, err = h.authService.GetUserByID(r.Context(), game.Player2ID)
			if err != nil {
				continue
			}
		}

		gameInfo := dto.ToGameShortInfo(game, player1, player2)
		gameInfos = append(gameInfos, gameInfo)
	}

	response := dto.AvailableGamesResponse{
		Games: gameInfos,
		Total: len(gameInfos),
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *GameHandler) GetUserCompletedGames(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	games, err := h.gameService.GetUserCompletedGames(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	var gameInfos []dto.GameShortInfo
	for _, game := range games {
		player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
		if err != nil {
			continue
		}

		var player2 *entities.User
		if game.Player2ID != "" && game.Player2ID != "AI" {
			player2, err = h.authService.GetUserByID(r.Context(), game.Player2ID)
			if err != nil {
				continue
			}
		}

		gameInfo := dto.ToGameShortInfo(game, player1, player2)
		gameInfos = append(gameInfos, gameInfo)
	}

	response := dto.AvailableGamesResponse{
		Games: gameInfos,
		Total: len(gameInfos),
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *GameHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	players, err := h.gameService.GetLeaderboard(r.Context(), limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	var leaderboardPlayers []dto.LeaderboardPlayer
	for _, player := range players {
		leaderboardPlayers = append(leaderboardPlayers, dto.LeaderboardPlayer{
			UUID:       player.UUID,
			Login:      player.Login,
			WinRatio:   player.WinRatio,
			Wins:       player.Wins,
			Losses:     player.Losses,
			Draws:      player.Draws,
			TotalGames: player.TotalGames,
		})
	}

	response := dto.LeaderboardResponse{
		Players: leaderboardPlayers,
		Total:   len(leaderboardPlayers),
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *GameHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case entities.ErrGameNotFound:
		h.sendErrorResponse(w, "Game not found", "GAME_NOT_FOUND", http.StatusNotFound)
	case entities.ErrUserNotFound:
		h.sendErrorResponse(w, "User not found", "USER_NOT_FOUND", http.StatusNotFound)
	case entities.ErrGameFull:
		h.sendErrorResponse(w, "Game is full", "GAME_FULL", http.StatusConflict)
	case entities.ErrWrongTurn:
		h.sendErrorResponse(w, "Not your turn", "WRONG_TURN", http.StatusBadRequest)
	case entities.ErrInvalidMove:
		h.sendErrorResponse(w, "Invalid move", "INVALID_MOVE", http.StatusBadRequest)
	case entities.ErrInvalidPosition:
		h.sendErrorResponse(w, "Position already occupied", "INVALID_POSITION", http.StatusBadRequest)
	case entities.ErrGameNotRunning:
		h.sendErrorResponse(w, "Game is not in progress", "GAME_NOT_RUNNING", http.StatusBadRequest)
	case entities.ErrForbidden:
		h.sendErrorResponse(w, "Forbidden", "FORBIDDEN", http.StatusForbidden)
	case entities.ErrUnauthorized:
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
	default:
		h.sendErrorResponse(w, "Internal server error", "INTERNAL_ERROR", http.StatusInternalServerError)
	}
}

func (h *GameHandler) sendErrorResponse(w http.ResponseWriter, message, code string, statusCode int) {
	middleware.SendErrorResponse(w, message, code, statusCode)
}

func (h *GameHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
