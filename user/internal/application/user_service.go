package application

import (
	"errors"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/ports"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	repository   ports.UserRepository
	tokenService ports.TokenService
}

func NewUserService(
	repository ports.UserRepository,
	tokenService ports.TokenService,
) *UserService {
	return &UserService{
		repository:   repository,
		tokenService: tokenService,
	}
}
