package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/domain"
	port "github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/ports"
)

type ProductService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) port.ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, in port.CreateProductInput) (*domain.Product, error) {
	now := time.Now().UTC()
	p := &domain.Product{
		ID:          uuid.NewString(),
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Get(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.List(ctx)
}

func (s *ProductService) Update(ctx context.Context, id string, in port.UpdateProductInput) (*domain.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		p.Name = in.Name
	}
	if in.Description != "" {
		p.Description = in.Description
	}
	if in.Price > 0 {
		p.Price = in.Price
	}
	if in.Stock != nil {
		p.Stock = *in.Stock
	}
	p.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) ReserveStock(ctx context.Context, id string, qty int) (*domain.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := p.ReserveStock(qty); err != nil {
		return nil, err
	}
	p.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
