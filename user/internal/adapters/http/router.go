package http

import (
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/http/middleware"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/ports"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	handler *Handler,
	tokenService ports.TokenService,
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)
	r.Post("/users", handler.CreateUser)
	r.Post("/users/login", handler.Login)

	r.With(middleware.Auth(tokenService)).Get(
		"/users/me",
		handler.GetMe,
	)

	return r
}
