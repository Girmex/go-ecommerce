package persistence

import (
	"context"

	"errors"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/domain"
	"gorm.io/gorm"
)

type CatalogRepository struct {
	db *gorm.DB
}

// DeleteProduct implements [ports.CatalogRepository].

func NewCatalogRepository(db *gorm.DB) *CatalogRepository {
	return &CatalogRepository{
		db: db,
	}
}


func (r *CatalogRepository) CreateCategory(ctx context.Context,category *domain.Category,) error {

	model := toCategoryModel(category)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Copy generated values back to the domain entity
	*category = *toCategoryDomain(model)

	return nil
}

func (repo *CatalogRepository) FindCategories(ctx context.Context,) ([]*domain.Category, error) {

	var models []models.CategoryModel

	err := repo.db.
		WithContext(ctx).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, 0, len(models))

	for _, model := range models {
		category := toCategoryDomain(&model)
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CatalogRepository) UpdateCategory(ctx context.Context,category *domain.Category,) (*domain.Category, error) {

	model := toCategoryModel(category)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	return toCategoryDomain(model), nil
}

func (repo *CatalogRepository) FindCategoryByID(
	ctx context.Context,
	id uint,
) (*domain.Category, error) {

	var model models.CategoryModel

	err := repo.db.WithContext(ctx).
		First(&model, id).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCategoryNotFound
		}

		return nil, err
	}

	return toCategoryDomain(&model), nil
}

func (repo *CatalogRepository) DeleteCategory(
	ctx context.Context,
	id uint,
) error {

	result := repo.db.
		WithContext(ctx).
		Delete(&models.CategoryModel{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (repo *CatalogRepository) CreateProduct(
	ctx context.Context,
	product *domain.Product,
) error {

	model := toProductModel(product)

	if err := repo.db.
		WithContext(ctx).
		Create(model).Error; err != nil {
		return err
	}

	*product = *toProductDomain(model)

	return nil
}

func (repo *CatalogRepository) UpdateProduct(
	ctx context.Context,
	product *domain.Product,
) (*domain.Product, error) {

	model := toProductModel(product)

	if err := repo.db.
		WithContext(ctx).
		Save(model).Error; err != nil {
		return nil, err
	}

	return toProductDomain(model), nil
}

func (repo *CatalogRepository) FindProductByID(
	ctx context.Context,
	id uint,
) (*domain.Product, error) {

	var model models.ProductModel

	err := repo.db.
		WithContext(ctx).
		First(&model, id).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}

		return nil, err
	}

	return toProductDomain(&model), nil
}

func (repo *CatalogRepository) FindProducts(
	ctx context.Context,
) ([]*domain.Product, error) {

	var models []models.ProductModel

	err := repo.db.
		WithContext(ctx).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	products := make([]*domain.Product, 0, len(models))

	for _, model := range models {
		products = append(products, toProductDomain(&model))
	}

	return products, nil
}

func (repo *CatalogRepository) FindSellerProducts(
	ctx context.Context,
	userID uint,
) ([]*domain.Product, error) {

	var models []models.ProductModel

	err := repo.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	products := make([]*domain.Product, 0, len(models))

	for _, model := range models {
		products = append(products, toProductDomain(&model))
	}

	return products, nil
}
func (r *CatalogRepository) DeleteProduct(ctx context.Context, product *domain.Product) error {
	panic("unimplemented")
}