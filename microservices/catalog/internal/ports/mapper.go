package persistence

import (
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/domain"
)

func toProductModel(product *domain.Product) *models.ProductModel {
	return &models.ProductModel{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		CategoryID:  product.CategoryID,
		ImageURL:    product.ImageURL,
		Price:       product.Price,
		UserID:      product.UserID,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func toProductDomain(model *models.ProductModel) *domain.Product {
	return &domain.Product{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CategoryID:  model.CategoryID,
		ImageURL:    model.ImageURL,
		Price:       model.Price,
		UserID:      model.UserID,
		Stock:       model.Stock,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func toCategoryModel(category *domain.Category) *models.CategoryModel {
	return &models.CategoryModel{
		ID:           category.ID,
		Name:         category.Name,
		ParentID:     category.ParentID,
		ImageURL:     category.ImageURL,
		DisplayOrder: category.DisplayOrder,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}

func toCategoryDomain(model *models.CategoryModel) *domain.Category {
	return &domain.Category{
		ID:           model.ID,
		Name:         model.Name,
		ParentID:     model.ParentID,
		ImageURL:     model.ImageURL,
		DisplayOrder: model.DisplayOrder,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}