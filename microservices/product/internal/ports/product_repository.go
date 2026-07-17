package ports

import "github.com/Girmex/go-ecommerce/microservices/product/internal/domain"

type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id uint) (*domain.Product, error)
	GetAll() ([]domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
}