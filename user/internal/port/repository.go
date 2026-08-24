package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
