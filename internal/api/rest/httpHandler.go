package rest

import (
	"github.com/Girmex/go-ecommerce/internal/helper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)
type RestHandler struct {
	App *fiber.App
	DB *gorm.DB
	Auth helper.Auth
}