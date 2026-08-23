package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/adapter/client"
	httpadapter "github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/adapter/http"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/adapter/persistence"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/application"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/config"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/database"
)

// @title           Order Service API
// @version         1.0
// @description     Order microservice for the eCommerce application.
// @host            localhost:8082
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token with the Bearer prefix.
//
// @tag.name orders
// @tag.description Order management endpoints
func main() {
	cfg := config.Load()

	ctx := context.Background()

	dbPool, err := database.NewPostgresPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(
		cfg.DatabaseURL,
		"migrations",
	); err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"Product service URL: %q",
		cfg.ProductServiceURL,
	)

	productClient := client.NewProductHTTPClient(
		cfg.ProductServiceURL,
	)

	paymentClient := client.NewPaymentHTTPClient(
		cfg.PaymentServiceURL,
	)

	orderRepository := persistence.NewOrderRepository(
		dbPool,
	)

	orderService := application.NewOrderService(
		orderRepository,
		productClient,
		paymentClient,
	)

	orderHandler := httpadapter.NewOrderHandler(
		orderService,
	)

	router := httpadapter.NewRouter(
		orderHandler,
		cfg.JWTSecret,
	)

	server := &nethttp.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf(
			"order service listening on :%s",
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

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("shutting down order service...")

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

	log.Println("order service stopped")
}
