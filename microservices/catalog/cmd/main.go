package main

import (
	"log"
	"net"

	grpcadapter "github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/grpc"
	grpcmiddleware "github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/grpc/middleware"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/database"
	"github.com/Girmex/go-ecommerce/microservices/catalog/proto"
	"github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	service := application.NewCatalogService(repository)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.AuthInterceptor(jwtManager),
		),
	)
	reflection.Register(server)
	handler := grpcadapter.NewHandler(service)
	proto.RegisterCatalogServiceServer(
		server,
		handler,
	)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(
		"%s started on :%s",
		cfg.AppName,
		cfg.GRPCPort,
	)
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
