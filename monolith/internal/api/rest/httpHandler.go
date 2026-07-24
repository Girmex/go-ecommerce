package rest

import (
	"github.com/Girmex/go-ecommerce/monolith/config"
	"github.com/Girmex/go-ecommerce/monolith/internal/helper"
	"github.com/Girmex/go-ecommerce/monolith/pkg/payment"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)
type RestHandler struct {
	App *fiber.App
	DB *gorm.DB
	Auth helper.Auth
	Config config.AppConfig
	PaymentClient payment.PaymentClient
}