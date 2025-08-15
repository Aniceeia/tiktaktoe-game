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

type UserHandler struct {
	authService *services.AuthService
}

func NewUserHandler(authService *services.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	currentUserID := middleware.GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	if userID == "" {
		userID = currentUserID
	}

	if currentUserID != userID {
		h.sendErrorResponse(w, "Forbidden", "FORBIDDEN", http.StatusForbidden)
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := dto.UserResponse{
		UUID:  user.UUID,
		Login: user.Login,
		Score: user.Score,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentUserID := middleware.GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}
	user, err := h.authService.GetUserByID(r.Context(), currentUserID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	response := dto.UserResponse{
		UUID:  user.UUID,
		Login: user.Login,
		Score: user.Score,
	}
	h.sendJSONResponse(w, response, http.StatusOK)
}

func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case entities.ErrUserNotFound:
		h.sendErrorResponse(w, "User not found", "USER_NOT_FOUND", http.StatusNotFound)
	case entities.ErrUnauthorized:
		h.sendErrorResponse(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
	case entities.ErrForbidden:
		h.sendErrorResponse(w, "Forbidden", "FORBIDDEN", http.StatusForbidden)
	default:
		h.sendErrorResponse(w, "Internal server error", "INTERNAL_ERROR", http.StatusInternalServerError)
	}
}

func (h *UserHandler) sendErrorResponse(w http.ResponseWriter, message, code string, statusCode int) {
	middleware.SendErrorResponse(w, message, code, statusCode)
}

func (h *UserHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
