package port

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, o *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Order, error)
	Update(ctx context.Context, o *domain.Order) error
}

type UserInfo struct {
	ID    string
	Name  string
	Email string
}

type UserClient interface {
	GetUser(ctx context.Context, userID string) (*UserInfo, error)
}

type ProductInfo struct {
	ID    string
	Name  string
	Price float64
	Stock int
}

type ProductClient interface {
	GetProduct(ctx context.Context, productID string) (*ProductInfo, error)
	ReserveStock(ctx context.Context, productID string, quantity int) error
}

type PaymentResult struct {
	ID     string
	Status string
}

type PaymentClient interface {
	Charge(ctx context.Context, orderID, userID string, amount float64, method string) (*PaymentResult, error)
}
