package http

import (
	"net/http"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/order/docs"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/adapter/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(h *OrderHandler, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"order-ms"}`))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))

		r.Post("/", h.Create)
		r.Get("/", h.ListByUser)
		r.Get("/{id}", h.Get)
		r.Post("/{id}/cancel", h.Cancel)
	})

	return r
}
