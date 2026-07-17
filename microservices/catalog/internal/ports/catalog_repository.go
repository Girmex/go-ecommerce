package ports

import "github.com/Girmex/go-ecommerce/microservices/catalog/internal/domain"

type CatalogRepository interface {

	// Category
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) (*domain.Category, error)
	DeleteCategory(id uint) error
	FindCategoryByID(id uint) (*domain.Category, error)
	FindCategories() ([]*domain.Category, error)

	// Product
	CreateProduct(product *domain.Product) error
	UpdateProduct(product *domain.Product) (*domain.Product, error)
	DeleteProduct(product *domain.Product) error
	FindProductByID(id uint) (*domain.Product, error)
	FindProducts() ([]*domain.Product, error)
	FindSellerProducts(userID uint) ([]*domain.Product, error)
}