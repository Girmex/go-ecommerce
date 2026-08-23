package port

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/domain"
)

type CreateOrderItemInput struct {
	ProductID string
	Quantity  int
}

type CreateOrderInput struct {
	UserID        uint
	Items         []CreateOrderItemInput
	PaymentMethod string
}

type OrderService interface {
	PlaceOrder(ctx context.Context, in CreateOrderInput) (*domain.Order, error)
	Get(ctx context.Context, id string) (*domain.Order, error)
	ListByUser(ctx context.Context, userID uint) ([]*domain.Order, error)
	Cancel(ctx context.Context, id string) (*domain.Order, error)
}
