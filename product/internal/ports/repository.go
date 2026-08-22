package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id string) error
}
