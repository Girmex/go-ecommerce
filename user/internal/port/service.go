package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"
)

type LoginResult struct {
	User  *domain.User
	Token string
}

type UserService interface {
	CreateUser(ctx context.Context, name string, email string, password string) (*domain.User, error)
	Login(ctx context.Context, email string, password string) (*LoginResult, error)
	GetUser(ctx context.Context, id uint) (*domain.User, error)
}