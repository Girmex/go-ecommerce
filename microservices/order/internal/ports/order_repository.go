package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error

	GetOrderByID(ctx context.Context, id uint) (*domain.Order, error)

	GetOrdersByUser(
		ctx context.Context,
		userID uint,
	) ([]domain.Order, error)

	UpdateOrder(ctx context.Context, order *domain.Order) error

	DeleteOrder(ctx context.Context, id uint) error
}
