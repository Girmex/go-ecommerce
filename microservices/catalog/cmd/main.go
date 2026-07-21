package main

import (
	"log"
	"net"

	"github.com/Girmex/go-ecommerce/microservices/catalog/api/proto"
	grpcadapter "github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/grpc"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/database"

	"google.golang.org/grpc"
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

	repository := persistence.NewCatalogRepository(db)

	service := application.NewCatalogService(repository)

	handler := grpcadapter.NewHandler(service)
	server := grpc.NewServer()
	proto.RegisterCatalogServiceServer(
		server,
		handler,
	)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Catalog Service Started")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
