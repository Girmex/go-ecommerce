package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/adapter/gateway"
	httpadapter "github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/adapter/http"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/adapter/persistence"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/application"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/config"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/database"

	_ "github.com/Girmex/go-ecommerce-app/chi-microservice/payment/docs"
)

// @title           Payment Service API
// @version         1.0
// @description     Payment microservice for the eCommerce application.
// @host            localhost:8084
// @BasePath        /
//
// @tag.name payments
// @tag.description Payment management endpoints
func main() {
	cfg := config.Load()

	ctx := context.Background()

	// ------------------------------------------------------------
	// Database
	// ------------------------------------------------------------

	dbPool, err := database.NewPostgresPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// ------------------------------------------------------------
	// Payment Gateway
	// ------------------------------------------------------------

	paymentGateway := gateway.NewMockGateway()

	// ------------------------------------------------------------
	// Persistence
	// ------------------------------------------------------------

	paymentRepository := persistence.NewPaymentRepository(
		dbPool,
	)

	// ------------------------------------------------------------
	// Application
	// ------------------------------------------------------------

	paymentService := application.NewPaymentService(
		paymentRepository,
		paymentGateway,
	)

	// ------------------------------------------------------------
	// HTTP
	// ------------------------------------------------------------

	paymentHandler := httpadapter.NewPaymentHandler(
		paymentService,
	)

	router := httpadapter.NewRouter(
		paymentHandler,
	)

	// ------------------------------------------------------------
	// HTTP Server
	// ------------------------------------------------------------

	server := &nethttp.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf(
			"payment service listening on :%s",
			cfg.HTTPPort,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != nethttp.ErrServerClosed {
			log.Fatalf(
				"failed to listen and serve: %v",
				err,
			)
		}
	}()

	// ------------------------------------------------------------
	// Graceful shutdown
	// ------------------------------------------------------------

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("shutting down payment service...")

	ctxShutdown, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf(
			"server forced to shutdown: %v",
			err,
		)
	}

	log.Println("payment service stopped")
}
