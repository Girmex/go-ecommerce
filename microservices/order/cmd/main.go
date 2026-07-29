package main

import (
	"log"
	"net"

	grpcadapter "github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/grpc"
	grpcmiddleware "github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/grpc/middleware"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/database"
	"github.com/Girmex/go-ecommerce/microservices/order/proto"
	"github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Load configuration
	cfg := config.Load()
	// Connect PostgreSQL
	pool, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(pool); err != nil {
		log.Fatal(err)
	}

	repository := persistence.NewOrderRepository(pool)

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	service := application.NewOrderService(repository)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.AuthInterceptor(jwtManager),
		),
	)
	reflection.Register(server)
	handler := grpcadapter.NewHandler(service)
	proto.RegisterOrderServiceServer(
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
