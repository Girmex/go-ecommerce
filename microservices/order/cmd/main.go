package main

import (
	"log"
	"net"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/cataloggrpc"
	grpcadapter "github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/grpc"
	grpcmiddleware "github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/grpc/middleware"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/database"
	"github.com/Girmex/go-ecommerce/microservices/order/proto"
	"github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	catalogConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer catalogConn.Close()

	repository := persistence.NewOrderRepository(pool)
	catalogClient := cataloggrpc.NewClient(catalogConn)

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	service := application.NewOrderService(repository, catalogClient)

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
