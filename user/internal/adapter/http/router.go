package http

import (
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapter/http/middleware"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/port"
	"github.com/go-chi/chi/v5"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/user/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	handler *Handler,
	tokenService ports.TokenService,
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// User routes
	r.Post("/users", handler.CreateUser)
	r.Post("/users/login", handler.Login)

	r.With(middleware.Auth(tokenService)).Get(
		"/users/me",
		handler.GetMe,
	)

	return r
}
