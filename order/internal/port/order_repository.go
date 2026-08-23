package port

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, o *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	ListByUser(ctx context.Context, userID uint) ([]*domain.Order, error)
	Update(ctx context.Context, o *domain.Order) error
}
