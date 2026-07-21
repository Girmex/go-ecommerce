package application

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application/dto"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/ports"
)

type CatalogService struct {
	repository ports.CatalogRepository
}

func NewCatalogService(repository ports.CatalogRepository) *CatalogService {
	return &CatalogService{
		repository: repository,
	}
}

func (s *CatalogService) CreateCategory(
	ctx context.Context,
	input dto.CreateCategoryInput,
) (*domain.Category, error) {

	category := &domain.Category{
		Name:         input.Name,
		ParentID:     input.ParentID,
		ImageURL:     input.ImageURL,
		DisplayOrder: input.DisplayOrder,
	}

	if err := s.repository.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	// category now contains the generated ID
	return category, nil
}

func (s *CatalogService) UpdateCategory(
	ctx context.Context,
	id uint,
	input dto.UpdateCategoryInput,
) (*domain.Category, error) {

	category, err := s.repository.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, domain.ErrCategoryNotFound
	}

	if input.Name != nil {
		category.Name = *input.Name
	}

	if input.ParentID != nil {
		category.ParentID = *input.ParentID
	}

	if input.ImageURL != nil {
		category.ImageURL = *input.ImageURL
	}

	if input.DisplayOrder != nil {
		category.DisplayOrder = *input.DisplayOrder
	}

	return s.repository.UpdateCategory(ctx, category)
}

func (s *CatalogService) DeleteCategory(
	ctx context.Context,
	id uint,
) error {

	if err := s.repository.DeleteCategory(ctx, id); err != nil {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (s *CatalogService) GetCategories(
	ctx context.Context,
) ([]*domain.Category, error) {

	return s.repository.FindCategories(ctx)
}

func (s *CatalogService) GetCategoryByID(
	ctx context.Context,
	id uint,
) (*domain.Category, error) {

	category, err := s.repository.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, domain.ErrCategoryNotFound
	}

	return category, nil
}

func (s *CatalogService) CreateProduct(
	ctx context.Context,
	input dto.CreateProductInput,
	userID uint,
) error {

	product := &domain.Product{
		Name:        input.Name,
		Description: input.Description,
		CategoryID:  input.CategoryID,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
		UserID:      userID,
		Stock:       uint(input.Stock),
	}

	return s.repository.CreateProduct(ctx, product)
}

func (s *CatalogService) UpdateProduct(
	ctx context.Context,
	id uint,
	input dto.UpdateProductInput,
	userID uint,
) (*domain.Product, error) {

	product, err := s.repository.FindProductByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}

	if product.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}

	if input.ImageURL != nil {
		product.ImageURL = *input.ImageURL
	}

	if input.Price != nil {
		product.Price = *input.Price
	}

	if input.Stock != nil {
		product.Stock = uint(*input.Stock)
	}

	return s.repository.UpdateProduct(ctx, product)
}

func (s *CatalogService) DeleteProduct(
	ctx context.Context,
	id uint,
	userID uint,
) error {

	product, err := s.repository.FindProductByID(ctx, id)
	if err != nil {
		return domain.ErrProductNotFound
	}

	if product.UserID != userID {
		return domain.ErrUnauthorized
	}

	return s.repository.DeleteProduct(ctx, product)
}

func (s *CatalogService) GetProducts(
	ctx context.Context,
) ([]*domain.Product, error) {

	return s.repository.FindProducts(ctx)
}

func (s *CatalogService) GetProductByID(
	ctx context.Context,
	id uint,
) (*domain.Product, error) {

	product, err := s.repository.FindProductByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}

	return product, nil
}

func (s *CatalogService) GetSellerProducts(
	ctx context.Context,
	userID uint,
) ([]*domain.Product, error) {

	return s.repository.FindSellerProducts(ctx, userID)
}

func (s *CatalogService) UpdateProductStock(
	ctx context.Context,
	id uint,
	input dto.UpdateStockInput,
	userID uint,
) (*domain.Product, error) {

	product, err := s.repository.FindProductByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}

	if product.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	product.Stock = uint(input.Stock)

	return s.repository.UpdateProduct(ctx, product)
}
