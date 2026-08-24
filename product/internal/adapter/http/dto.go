package http

import (
	"time"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/domain"
)

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=150" example:"Wireless Mouse"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0" example:"29.99"`
	Stock       int     `json:"stock" validate:"required,gte=0" example:"100"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=2,max=150"`
	Description string  `json:"description" validate:"omitempty,max=1000"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Stock       *int    `json:"stock" validate:"omitempty,gte=0"`
}

type ReserveStockRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0" example:"2"`
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ToProductResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}