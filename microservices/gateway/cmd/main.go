package main

import (
	"log"
	"net/http"

	authgrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/auth"
	cataloggrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/catalog"
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

	authHandler := handlers.NewAuthHandler(authClient)
	catalogHandler := handlers.NewCatalogHandler(catalogClient)
	r := chi.NewRouter()

	routes.RegisterRoutes(
		r,
		authHandler,
		catalogHandler,
		jwtManager,
	)
	log.Printf("%s started on :%s", cfg.AppName, cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
