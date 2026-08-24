package main

import (
	"context"
	"log"
	"net/http"

	httpHandler "github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/adapter/http"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/adapter/persistence"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/application"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/config"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/database"
)

// @title Product Service API
// @version 1.0
// @description Product microservice for the eCommerce application.
// @host localhost:8081
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token with the Bearer prefix.
func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(
		cfg.DatabaseURL,
		"migrations",
	); err != nil {
		log.Fatal(err)
	}

	// Persistence adapter
	productRepository := persistence.NewProductRepository(db)

	// Application service
	productService := application.NewProductService(
		productRepository,
	)

	// HTTP adapter
	handler := httpHandler.NewHandler(productService)

	// HTTP router
	router := httpHandler.NewRouter(
		handler,
		cfg.JwtSecret,
	)

	log.Printf(
		"product service listening on :%s",
		cfg.HTTPPort,
	)

	if err := http.ListenAndServe(
		":"+cfg.HTTPPort,
		router,
	); err != nil {
		log.Fatal(err)
	}
}