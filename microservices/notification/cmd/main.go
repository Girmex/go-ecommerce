package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/authgrpc"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/email"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/kafka"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/sms"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/config"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting %s in %s mode...\n", cfg.AppName, cfg.AppEnv)

	// Connect to Auth gRPC service for resolving user profiles
	authConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to Auth gRPC service: %v\n", err)
	}
	defer authConn.Close()

	authClient := authgrpc.NewClient(authConn)

	// Initialize Email Sender adapter
	var emailSender ports.EmailSender
	if cfg.EmailProvider == "smtp" {
		log.Printf("Using SMTP Email Sender (%s:%s)\n", cfg.SMTPHost, cfg.SMTPPort)
		emailSender = email.NewSMTPEmailSender(
			cfg.SMTPHost,
			cfg.SMTPPort,
			cfg.SMTPUsername,
			cfg.SMTPPassword,
			cfg.SMTPFrom,
		)
	} else {
		log.Println("Using Logger Email Sender (Console)")
		emailSender = email.NewLoggerEmailSender()
	}

	// Initialize SMS Sender adapter
	var smsSender ports.SMSSender
	if cfg.SMSProvider == "twilio" {
		log.Println("Using Twilio SMS Sender")
		smsSender = sms.NewTwilioSMSSender(
			cfg.TwilioAccountSID,
			cfg.TwilioAuthToken,
			cfg.TwilioFromPhone,
		)
	} else {
		log.Println("Using Logger SMS Sender (Console)")
		smsSender = sms.NewLoggerSMSSender()
	}

	// Initialize Application Service
	notificationService := application.NewNotificationService(emailSender, smsSender, authClient)

	// Context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize and start Kafka Listener
	listener := kafka.NewEventListener(
		cfg.KAFKABrokers,
		cfg.ConsumerGroupID,
		notificationService,
	)

	log.Printf("Connecting Kafka consumers to brokers: %v (group: %s)\n", cfg.KAFKABrokers, cfg.ConsumerGroupID)

	if err := listener.Start(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Notification service error: %v\n", err)
	}

	log.Println("Notification service shut down gracefully.")
}
