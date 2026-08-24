package http

import (
	"net/http"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/product/docs"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/adapter/http/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(productHandler *Handler, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
	))

	r.Route("/products", func(r chi.Router) {
		// Public endpoints
		r.Get("/", productHandler.List)
		r.Get("/{id}", productHandler.Get)

		// Authenticated endpoints
		r.With(middleware.Auth(jwtSecret)).Post("/", productHandler.Create)
		r.With(middleware.Auth(jwtSecret)).Put("/{id}", productHandler.Update)
		r.With(middleware.Auth(jwtSecret)).Delete("/{id}", productHandler.Delete)

		// Internal service endpoint
		r.Post("/{id}/reserve", productHandler.ReserveStock)
	})

	return r
}