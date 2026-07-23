package ports
import (
	"context"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
)

type AuthRepository interface {

    CreateUser(
        ctx context.Context,
        user *domain.User,
    ) error

    GetUserByID(
        ctx context.Context,
        id uint,
    ) (*domain.User, error)

    GetUserByEmail(
        ctx context.Context,
        email string,
    ) (*domain.User, error)
}