package main

import (
	"log"
	"net/http"

	authgrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/auth"
	cataloggrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/catalog"
	ordergrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/order"
	paymentgrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/payment"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/http/handlers"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/routes"
	jwtpkg "github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.Load()
	jwtManager := jwtpkg.NewJWTManager(cfg.JWTSecret)

	authClient, err := authgrpc.New(cfg.GRPCAuthHost)
	if err != nil {
		log.Fatal(err)
	}

	catalogClient, err := cataloggrpc.New(cfg.GRPCCatalogHost)
	if err != nil {
		log.Fatal(err)
	}
	orderClient, err := ordergrpc.New(cfg.GRPCOrderHost)

	if err != nil {
		log.Fatal(err)
	}

	paymentClient, err := paymentgrpc.New(cfg.GRPCPaymentHost)
	if err != nil {
		log.Fatal(err)
	}

	authHandler := handlers.NewAuthHandler(authClient)
	catalogHandler := handlers.NewCatalogHandler(catalogClient)
	orderHandler := handlers.NewOrderHandler(orderClient)
	paymentHandler := handlers.NewPaymentHandler(paymentClient)

	r := chi.NewRouter()

	routes.RegisterRoutes(
		r,
		authHandler,
		catalogHandler,
		orderHandler,
		paymentHandler,
		jwtManager,
	)
	log.Printf("%s started on :%s", cfg.AppName, cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
