package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/adapter/http/middleware"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/payment/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(h *PaymentHandler, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(20 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(
			[]byte(`{"status":"ok","service":"payment-ms"}`),
		)
	})

	r.Get(
		"/swagger/*",
		httpSwagger.Handler(
			httpSwagger.URL(
				"http://localhost:8083/swagger/doc.json",
			),
		),
	)

	r.Route("/payments", func(r chi.Router) {
		// JWT authentication
		r.Use(middleware.Auth(jwtSecret))

		r.Post("/", h.Charge)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Post("/{id}/refund", h.Refund)
	})

	return r
}