package http

import (
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	userService *application.UserService
	validator   *validator.Validate
}

func NewHandler(userService *application.UserService) *Handler {
	return &Handler{
		userService: userService,
		validator:   validator.New(),
	}
}