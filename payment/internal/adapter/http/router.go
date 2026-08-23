package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/payment/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(h *PaymentHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(20 * time.Second))

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
				"http://localhost:8084/swagger/doc.json",
			),
		),
	)

	r.Route("/payments", func(r chi.Router) {
		r.Post("/", h.Charge)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Post("/{id}/refund", h.Refund)
	})

	return r
}
