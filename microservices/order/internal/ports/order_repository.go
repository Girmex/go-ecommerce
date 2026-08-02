package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error

	GetOrder(ctx context.Context, orderID uint) (*domain.Order, error)

	GetOrdersByUser(ctx context.Context, userID uint) ([]domain.Order, error)

	UpdateOrder(ctx context.Context, order *domain.Order) error

	DeleteOrder(ctx context.Context, orderID uint) error
	UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) (*domain.Order, error)
}
