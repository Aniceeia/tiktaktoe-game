package api

import (
	"encoding/json"
	"game/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GameHandler struct {
	service     service.GameService
	ai          *service.AI
	authService service.AuthService // Используем интерфейс вместо конкретной реализации
}

func NewGameHandler(
	service service.GameService,
	ai *service.AI,
	authService service.AuthService, // Используем интерфейс
) *GameHandler {
	return &GameHandler{
		service:     service,
		ai:          ai,
		authService: authService,
	}
}

// getCurrentUser extracts the user UUID from context (set by UserAuthenticator)
func getCurrentUser(r *http.Request) string {
	ctx := r.Context()
	uuid, ok := ctx.Value(userUUIDKey).(string)
	if !ok {
		return ""
	}
	return uuid
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	userUUID := getCurrentUser(r)
	if userUUID == "" {
		respondError(w, "User authentication required", http.StatusUnauthorized)
		return
	}
	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Mode != "pvp" && req.Mode != "pve" {
		respondError(w, "Invalid mode. Must be 'pvp' or 'pve'", http.StatusBadRequest)
		return
	}
	newGame, err := h.service.CreateGame(userUUID, req.Mode)
	if err != nil {
		respondError(w, "Failed to create game", http.StatusInternalServerError)
		return
	}
	respondJSON(w, ToResponse(newGame), http.StatusCreated)
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")
	userUUID := getCurrentUser(r)
	if userUUID == "" {
		respondError(w, "User authentication required", http.StatusUnauthorized)
		return
	}
	err := h.service.JoinGame(gameID, userUUID)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, map[string]string{"message": "Successfully joined game"}, http.StatusOK)
}

func (h *GameHandler) HandleMove(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")
	userUUID := getCurrentUser(r)
	if userUUID == "" {
		respondError(w, "User authentication required", http.StatusUnauthorized)
		return
	}
	var req UpdateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request handle move", http.StatusBadRequest)
		return
	}
	updatedGame, err := h.service.ProcessMove(gameID, userUUID, req.Board)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, ToResponse(updatedGame), http.StatusAccepted)
}

func (h *GameHandler) GetGameState(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")

	game, err := h.service.GetGameState(gameID)
	if err != nil {
		respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(w, ToResponse(game), http.StatusOK)
}

func (h *GameHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		respondError(w, "Login and password required", http.StatusBadRequest)
		return
	}
	if err := h.authService.Register(service.SignUpRequest{
		Login:    req.Login,
		Password: req.Password,
	}); err != nil {
		if err == service.ErrUserExists {
			respondError(w, "User already exists", http.StatusBadRequest)
			return
		}
		respondError(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	respondJSON(w, map[string]string{"message": "User registered successfully"}, http.StatusCreated)
}

func respondJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, ErrorResponse{Error: message}, statusCode)
}

func (h *GameHandler) Login(w http.ResponseWriter, r *http.Request) {
	login, password, ok := r.BasicAuth()
	if !ok {
		respondError(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	uuid, err := h.authService.Authenticate(login, password)
	if err != nil {
		respondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	respondJSON(w, map[string]string{"uuid": uuid}, http.StatusOK)
}

//last of 31

func (h *GameHandler) GetAvailableGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.service.GetAvailableGames()
	if err != nil {
		respondError(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	response := make([]GameResponse, len(games))
	for i, game := range games {
		response[i] = ToResponse(game)
	}
	respondJSON(w, response, http.StatusOK)
}

func (h *GameHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	user, err := h.authService.GetUserByUUID(uuid)
	if err != nil {
		respondError(w, "User not found", http.StatusNotFound)
		return
	}

	respondJSON(w, UserResponse{
		UUID:     user.UUID,
		Username: user.Username,
		Score:    user.Score,
	}, http.StatusOK)
}
