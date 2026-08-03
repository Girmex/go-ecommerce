package main

import (
	"log"
	"net/http"

	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/config"
	grpcclient "github.com/Girmex/go-ecommerce/microservices/gateway/internal/grpc"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/handlers"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/routes"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.Load()

	clients, err := grpcclient.NewClients(cfg.GRPCAuthHost)
	if err != nil {
		log.Fatal(err)
	}

	authHandler := handlers.NewAuthHandler(clients.Auth)

	r := chi.NewRouter()

	routes.RegisterRoutes(r, authHandler)

	log.Printf("%s started on :%s", cfg.AppName, cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
