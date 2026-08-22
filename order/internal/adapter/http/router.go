package http

import (
	"net/http"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/order/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(h *OrderHandler) http.Handler {
	r := chi.NewRouter()

			r.Use(middleware.RequestID)
			r.Use(middleware.Logger)
			r.Use(middleware.Recoverer)
			r.Use(middleware.Timeout(20))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"order-ms"}`))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.ListByUser)
		r.Get("/{id}", h.Get)
		r.Post("/{id}/cancel", h.Cancel)
	})

	return r
}
