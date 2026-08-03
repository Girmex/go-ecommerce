package routes

import (
	"net/http"

	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, authHandler *handlers.AuthHandler) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gateway running"))
	})

	// auth routes
	r.Route("/auth", func(r chi.Router) {

		// next we'll implement these
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})
}