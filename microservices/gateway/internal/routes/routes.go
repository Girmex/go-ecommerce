package routes

import (
	"net/http"

	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/http/handlers"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/middleware"
	jwtpkg "github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r *chi.Mux,
	authHandler *handlers.AuthHandler,
	catalogHandler *handlers.CatalogHandler,
	jwtManager *jwtpkg.JWTManager,
) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gateway running"))
	})

	// auth routes
	r.Route("/auth", func(r chi.Router) {

		// next we'll implement these
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/catalog", func(r chi.Router) {

		r.Get("/products", catalogHandler.ListProducts)

		r.Group(func(r chi.Router) {

			r.Use(middleware.JWTMiddleware(jwtManager))

			r.Post("/products", catalogHandler.CreateProduct)

		})
	})
}
