package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/auth"
	httpadapter "github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/http"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/persistence"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/config"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/database"
)

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

	handler := httpadapter.NewHandler(userService)
	router := httpadapter.NewRouter(handler)

	log.Printf("user service listening on :%s", cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		log.Fatal(err)
	}
}
