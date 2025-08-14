package handlers

import (
	"encoding/json"
	"net/http"

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

// CreateGame создает новую игру
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	// Валидация режима игры
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

	// Получаем информацию о пользователях для ответа
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

// JoinGame присоединяет игрока к существующей игре
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

	// Получаем обновленное состояние игры
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

// MakeMove выполняет ход в игре
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

	// Получаем информацию о пользователях для ответа
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

// GetGameState возвращает текущее состояние игры
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

	// Проверяем, что пользователь является участником игры
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

// GetAvailableGames возвращает список доступных игр
func (h *GameHandler) GetAvailableGames(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	games, err := h.gameService.GetAvailableGames(r.Context())
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	var gameInfos []dto.GameShortInfo
	for _, game := range games {
		player1, err := h.authService.GetUserByID(r.Context(), game.Player1ID)
		if err != nil {
			continue // Пропускаем игры с проблемными пользователями
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

// GetUserGames возвращает игры пользователя
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

// handleServiceError обрабатывает ошибки сервиса и отправляет соответствующий HTTP статус
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

// sendErrorResponse отправляет ответ с ошибкой
func (h *GameHandler) sendErrorResponse(w http.ResponseWriter, message, code string, statusCode int) {
	response := dto.ErrorResponse{
		Error: message,
		Code:  code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// sendJSONResponse отправляет JSON ответ
func (h *GameHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
