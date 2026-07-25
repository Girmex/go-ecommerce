package main

import (
	"log"
	"net"

	grpcadapter "github.com/Girmex/go-ecommerce/microservices/auth/internal/adapters/grpc"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/database"
	"github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"github.com/Girmex/go-ecommerce/microservices/auth/proto"
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
		&models.UserModel{},
	); err != nil {
		log.Fatal(err)
	}

	repository := persistence.NewAuthRepository(db)

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	service := application.NewAuthService(
		repository,
		jwtManager,
	)
	handler := grpcadapter.NewHandler(service)

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(
		server,
		handler,
	)
	reflection.Register(server)

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
