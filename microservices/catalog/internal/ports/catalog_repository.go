package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/domain"
)

type CatalogRepository interface {

	// Category
	CreateCategory(ctx context.Context, category *domain.Category) error
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	FindCategoryByID(ctx context.Context, id uint) (*domain.Category, error)
	FindCategories(ctx context.Context) ([]*domain.Category, error)

	// Product
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, product *domain.Product) error
	FindProductByID(ctx context.Context, id uint) (*domain.Product, error)
	FindProducts(ctx context.Context) ([]*domain.Product, error)
	FindSellerProducts(ctx context.Context, userID uint) ([]*domain.Product, error)
}