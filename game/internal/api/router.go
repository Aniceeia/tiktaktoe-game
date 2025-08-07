package api

import (
	"github.com/go-chi/chi/v5"
)

type Router chi.Router

func NewRouter(gh *GameHandler) Router {
	r := chi.NewRouter()
	r.Post("/register", gh.Register)
	r.Post("/login", gh.Login)
	r.Group(func(protected chi.Router) {
		protected.Use(gh.UserAuthenticator)
		protected.Get("/games/available", gh.GetAvailableGames)
		protected.Get("/users/{uuid}", gh.GetUser)
		protected.Route("/games", func(r chi.Router) {
			r.Post("/", gh.CreateGame)
			r.Get("/{id}", gh.GetGameState)
			r.Post("/{id}/join", gh.JoinGame)
			r.Post("/{id}/move", gh.HandleMove)
		})
	})
	return r
}
