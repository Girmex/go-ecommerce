package routers

import (
	"net/http"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/product/docs"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/adapters/http/handlers"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(productHandler *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
	))
	r.Route("/products", func(r chi.Router) {
		r.Post("/", productHandler.Create)
		r.Get("/", productHandler.List)

		r.Get("/{id}", productHandler.Get)
		r.Put("/{id}", productHandler.Update)
		r.Delete("/{id}", productHandler.Delete)

		r.Post("/{id}/reserve", productHandler.ReserveStock)
	})

	return r
}
