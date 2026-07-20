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

	log.Printf("Host=%q", cfg.PostgresHost)
	log.Printf("Port=%q", cfg.PostgresPort)
	log.Printf("User=%q", cfg.PostgresUser)
	log.Printf("Password=%q", cfg.PostgresPassword)
	log.Printf("Database=%q", cfg.PostgresDatabase)

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

	repository := persistence.NewCatalogRepository(db)

	service := application.NewCatalogService(repository)

	_ = service

	log.Println("Catalog Service Started")
}
