package main

import (
	"context"
	"log"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/config"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/database"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(
		cfg.DatabaseURL,
		"migrations",
	); err != nil {
		log.Fatal(err)
	}

	userRepository := persistence.NewUserRepository(db)

	userService := application.NewUserService(userRepository)

	_ = userService

	log.Println("user service started")
}
