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
	smsadapter "github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/sms"
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
	// SMS adapter
	// --------------------------------------------------

	//smsSender := smsadapter.NewLoggingSMSSender()
	smsSender := smsadapter.NewTwilioSMSSender(
		cfg.TwilioAccountSID,
		cfg.TwilioAuthToken,
		cfg.TwilioFromPhone,
	)
	// --------------------------------------------------
	// Application service
	// --------------------------------------------------

	notificationService := application.NewNotificationService(
		userClient,
		emailSender,
		smsSender,
	)

	// --------------------------------------------------
	// Payment completed Kafka consumer
	// --------------------------------------------------

	paymentKafkaConsumer := kafkapkg.NewConsumer(
		[]string{cfg.KAFKABrokers},
		kafkapkg.TopicPaymentCompleted,
		"notification-payment-service",
	)

	defer paymentKafkaConsumer.Close()

	paymentConsumer := kafkaadapter.NewPaymentCompletedConsumer(
		paymentKafkaConsumer,
		notificationService,
	)

	// --------------------------------------------------
	// Phone verification Kafka consumer
	// --------------------------------------------------

	phoneKafkaConsumer := kafkapkg.NewConsumer(
		[]string{cfg.KAFKABrokers},
		kafkapkg.TopicUserPhoneVerification,
		"notification-phone-verification-service",
	)

	defer phoneKafkaConsumer.Close()

	phoneConsumer := kafkaadapter.NewPhoneVerificationConsumer(
		phoneKafkaConsumer,
		notificationService,
	)

	// --------------------------------------------------
	// Start
	// --------------------------------------------------

	log.Printf(
		"%s started",
		cfg.AppName,
	)

	errCh := make(chan error, 2)

	go func() {
		errCh <- paymentConsumer.Start(ctx)
	}()

	go func() {
		errCh <- phoneConsumer.Start(ctx)
	}()

	select {
	case err := <-errCh:
		if ctx.Err() != nil {
			log.Println("Notification Service shutting down")
			return
		}

		log.Fatalf(
			"notification consumer: %v",
			err,
		)

	case <-ctx.Done():
		log.Println("Notification Service shutting down")
	}
}
