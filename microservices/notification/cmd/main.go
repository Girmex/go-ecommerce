package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	authgrpc "github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/authgrpc"
	emailadapter "github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/email"
	kafkaadapter "github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/kafka"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/config"
	kafkapkg "github.com/Girmex/go-ecommerce/microservices/pkg/kafka"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	// Context cancelled when the process receives SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// --------------------------------------------------
	// Auth gRPC connection
	// --------------------------------------------------

	

	authConn, err := grpc.NewClient(
		cfg.AuthGRPCAddress,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer authConn.Close()

	userClient := authgrpc.NewClient(authConn)

	// --------------------------------------------------
	// Email adapter
	// --------------------------------------------------

	emailSender := emailadapter.NewSMTPSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

	// --------------------------------------------------
	// Application service
	// --------------------------------------------------

	notificationService := application.NewNotificationService(
		userClient,
		emailSender,
	)

	// --------------------------------------------------
	// Kafka consumer
	// --------------------------------------------------

	consumer := kafkapkg.NewConsumer(
		[]string{cfg.KAFKABrokers},
		kafkapkg.TopicPaymentCompleted,
		"notification-service",
	)

	defer consumer.Close()

	paymentConsumer := kafkaadapter.NewPaymentCompletedConsumer(
		consumer,
		notificationService,
	)

	// --------------------------------------------------
	// Start
	// --------------------------------------------------

	log.Printf(
		"%s started",
		cfg.AppName,
	)

	if err := paymentConsumer.Start(ctx); err != nil {
		if ctx.Err() != nil {
			log.Println("Notification Service shutting down")
			return
		}

		log.Fatalf(
			"notification consumer: %v",
			err,
		)
	}
}
