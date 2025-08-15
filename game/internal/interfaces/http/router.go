package http

import (
	"encoding/json"
	"game/internal/domain/services"
	"game/internal/interfaces/http/handlers"
	"game/internal/interfaces/http/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	router      *mux.Router
	gameHandler *handlers.GameHandler
	authHandler *handlers.AuthHandler
	userHandler *handlers.UserHandler
	authService *services.AuthService
}

func NewRouter(
	gameService *services.GameService,
	authService *services.AuthService,
) *Router {
	router := mux.NewRouter()

	gameHandler := handlers.NewGameHandler(gameService, authService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	r := &Router{
		router:      router,
		gameHandler: gameHandler,
		authHandler: authHandler,
		userHandler: userHandler,
		authService: authService,
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {

	r.router.Use(middleware.CORSMiddleware())
	r.router.Use(middleware.RecoveryMiddleware())
	r.router.Use(middleware.LoggingMiddleware())

	api := r.router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/auth/register", r.authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", r.authHandler.Login).Methods("POST")

	api.Handle("/games", middleware.AuthMiddleware(r.authService)(middleware.ValidationMiddleware()(http.HandlerFunc(r.gameHandler.CreateGame)))).Methods("POST")
	api.Handle("/games", middleware.AuthMiddleware(r.authService)(http.HandlerFunc(r.gameHandler.GetAvailableGames))).Methods("GET")
	api.Handle("/games/my", middleware.AuthMiddleware(r.authService)(http.HandlerFunc(r.gameHandler.GetUserGames))).Methods("GET")
	api.Handle("/games/{gameId}", middleware.AuthMiddleware(r.authService)(http.HandlerFunc(r.gameHandler.GetGameState))).Methods("GET")
	api.Handle("/games/{gameId}/join", middleware.AuthMiddleware(r.authService)(http.HandlerFunc(r.gameHandler.JoinGame))).Methods("POST")
	api.Handle("/games/{gameId}/move", middleware.AuthMiddleware(r.authService)(middleware.ValidationMiddleware()(http.HandlerFunc(r.gameHandler.MakeMove)))).Methods("POST")

	api.Handle("/users/{userId}", middleware.AuthMiddleware(r.authService)(http.HandlerFunc(r.userHandler.GetUser))).Methods("GET")

	r.router.HandleFunc("/health", r.healthCheck).Methods("GET")
}

func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "tic-tac-toe-api",
	})
}

func (r *Router) GetRouter() *mux.Router {
	return r.router
}
