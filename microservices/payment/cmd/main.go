package main

import (
	"log"
	"net"

	grpcadapter "github.com/Girmex/go-ecommerce/microservices/payment/internal/adapters/grpc"
	grpcmiddleware "github.com/Girmex/go-ecommerce/microservices/payment/internal/adapters/grpc/middleware"
	ordergrpc "github.com/Girmex/go-ecommerce/microservices/payment/internal/adapters/ordergrpc"
	"github.com/Girmex/go-ecommerce/microservices/payment/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce/microservices/payment/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/payment/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/payment/internal/database"
	"github.com/Girmex/go-ecommerce/microservices/payment/proto"
	"github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"github.com/Girmex/go-ecommerce/microservices/pkg/kafka"
	"google.golang.org/grpc/credentials/insecure"

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
	orderConn, err := grpc.NewClient(
		"localhost:50053",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer orderConn.Close()

	repository := persistence.NewPaymentRepository(pool)

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	orderClient := ordergrpc.NewClient(orderConn)

	kafkaProducer := kafka.NewProducer(
		[]string{cfg.KAFKABrokers},
	)
	defer kafkaProducer.Close()

	service := application.NewPaymentService(repository, orderClient, kafkaProducer)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.AuthInterceptor(jwtManager),
		),
	)

	reflection.Register(server)

	handler := grpcadapter.NewHandler(service)

	proto.RegisterPaymentServiceServer(
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
