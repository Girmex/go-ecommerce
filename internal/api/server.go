package api

import (
	"log"
	"net/http"

	"github.com/Girmex/go-ecommerce/config"
	"github.com/Girmex/go-ecommerce/internal/api/rest"
	"github.com/Girmex/go-ecommerce/internal/api/rest/handlers"
	"github.com/Girmex/go-ecommerce/internal/domain"
	"github.com/Girmex/go-ecommerce/internal/helper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig){
	app := fiber.New()	

	db, err:= gorm.Open(postgres.Open(config.Dsn),&gorm.Config{})

	if err!=nil{
		log.Fatalf("Database connection error %v\n",err)
	}

	log.Println("database connected!\n",db)

	err= db.AutoMigrate(&domain.User{},&domain.BankAccount{},&domain.Category{},&domain.Product{})
	if err != nil {
		log.Fatalf("error on runing migration %v", err.Error())
	}

	log.Println("migration was successful")

	auth := helper.SetupAuth(config.AppSecret)

	rh:= &rest.RestHandler{
		App: app,
		DB: db,
		Auth: auth,
		Config: config,


	}

	log.Printf("config DSN: %v",config.Dsn)

	setupRoutes(rh)

	app.Get("/health",HealthCheck)
	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.RestHandler){
	handlers.SetupUserRoutes(rh)
}

func HealthCheck(ctx * fiber.Ctx) error{
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"Message":"I am breathing",
	})
}