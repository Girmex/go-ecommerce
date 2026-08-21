package http

import "github.com/go-chi/chi/v5"

func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)
	r.Post("/users", handler.CreateUser)
	r.Post("/users/login", handler.Login)
	return r
}
