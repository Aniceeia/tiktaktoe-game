package handlers

import (
	"encoding/json"
	"net/http"

	"game/internal/domain/entities"
	"game/internal/domain/services"
	"game/internal/interfaces/http/dto"
	"game/internal/interfaces/http/middleware"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	// Валидация запроса
	if err := validate.Struct(req); err != nil {
		h.sendErrorResponse(w, "Validation failed", "VALIDATION_FAILED", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := dto.UserResponse{
		UUID:  user.UUID,
		Login: user.Login,
		Score: user.Score,
	}

	h.sendJSONResponse(w, response, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	login, password, ok := r.BasicAuth()
	if !ok {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	// Аутентифицируем пользователя
	user, err := h.authService.Login(r.Context(), login, password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := dto.LoginResponse{
		UUID: user.UUID,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

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

func (h *AuthHandler) sendErrorResponse(w http.ResponseWriter, message, code string, statusCode int) {
	middleware.SendErrorResponse(w, message, code, statusCode)
}

func (h *AuthHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
