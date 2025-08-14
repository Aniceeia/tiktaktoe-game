package handlers

import (
	"encoding/json"
	"net/http"

	"game/internal/domain/entities"
	"game/internal/domain/services"
	"game/internal/interfaces/http/dto"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register регистрирует нового пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := dto.UserResponse{
		UUID:     user.UUID,
		Username: user.Username,
		Score:    user.Score,
	}

	h.sendJSONResponse(w, response, http.StatusCreated)
}

// Login аутентифицирует пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Получаем Basic Auth данные
	username, password, ok := r.BasicAuth()
	if !ok {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	// Аутентифицируем пользователя
	user, err := h.authService.Login(r.Context(), username, password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := dto.UserResponse{
		UUID:     user.UUID,
		Username: user.Username,
		Score:    user.Score,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// handleServiceError обрабатывает ошибки сервиса
func (h *AuthHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case entities.ErrUserExists:
		h.sendErrorResponse(w, "User already exists", "USER_EXISTS", http.StatusConflict)
	case entities.ErrInvalidCredentials:
		h.sendErrorResponse(w, "Invalid credentials", "INVALID_CREDENTIALS", http.StatusUnauthorized)
	case entities.ErrUserNotFound:
		h.sendErrorResponse(w, "User not found", "USER_NOT_FOUND", http.StatusNotFound)
	default:
		h.sendErrorResponse(w, "Internal server error", "INTERNAL_ERROR", http.StatusInternalServerError)
	}
}

// sendErrorResponse отправляет ответ с ошибкой
func (h *AuthHandler) sendErrorResponse(w http.ResponseWriter, message, code string, statusCode int) {
	response := dto.ErrorResponse{
		Error: message,
		Code:  code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// sendJSONResponse отправляет JSON ответ
func (h *AuthHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
