package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapter/auth"
	httpahandler "github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapter/http"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapter/persistence"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/config"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/database"
)

// @title User Service API
// @version 1.0
// @description User microservice for the eCommerce application.
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
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

	tokenService := auth.NewJWTService(cfg.JWTSecret)

	userService := application.NewUserService(
		userRepository,
		tokenService,
	)

	handler := httpahandler.NewHandler(userService)
	router := httpahandler.NewRouter(handler, tokenService)

	log.Printf("user service listening on :%s", cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		log.Fatal(err)
	}
}
