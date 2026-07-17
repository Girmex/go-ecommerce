package main

import (
	"log"

	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/database"
)

func main() {

	// Load configuration
	cfg := config.Load()

	// Connect PostgreSQL
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate persistence models
	if err := db.AutoMigrate(
		&models.CategoryModel{},
		&models.ProductModel{},
	); err != nil {
		log.Fatal(err)
	}

	// Repository
	repository := persistence.NewCatalogRepository(db)

	// Application Service
	service := application.NewCatalogService(repository)

	// Avoid unused variable until gRPC server is added
	_ = service

	log.Println("Catalog Service Started")
}