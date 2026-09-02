package port

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/domain"
)

type CreateProductInput struct {
	UserID      uint
	Name        string
	Description string
	Price       float64
	Stock       int
}

type UpdateProductInput struct {
	UserID      uint
	Name        string
	Description string
	Price       float64
	Stock       *int
}

type ProductService interface {
	Create(ctx context.Context, in CreateProductInput) (*domain.Product, error)
	Get(ctx context.Context, id string) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Update(ctx context.Context, id string, in UpdateProductInput) (*domain.Product, error)
	Delete(ctx context.Context, id string, userID uint) error
	ReserveStock(ctx context.Context, id string, qty int) (*domain.Product, error)
}
