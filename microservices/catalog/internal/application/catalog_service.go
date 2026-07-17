package application

import (
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

func (s *CatalogService) CreateCategory(input dto.CreateCategoryInput) error {
	return s.repository.CreateCategory(&domain.Category{
		Name:         input.Name,
		ParentID:     input.ParentID,
		ImageURL:     input.ImageURL,
		DisplayOrder: input.DisplayOrder,
	})
}

func (s *CatalogService) UpdateCategory(id uint, input dto.UpdateCategoryInput) (*domain.Category, error) {

	category, err := s.repository.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("category does not exist")
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

	return s.repository.UpdateCategory(category)
}
func (s *CatalogService) DeleteCategory(id uint) error {

	err := s.repository.DeleteCategory(id)
	if err != nil {
		return errors.New("category does not exist")
	}

	return nil
}

func (s *CatalogService) GetCategories() ([]*domain.Category, error) {

	categories, err := s.repository.FindCategories()
	if err != nil {
		return nil, errors.New("categories do not exist")
	}

	return categories, nil
}

func (s *CatalogService) GetCategoryByID(id uint) (*domain.Category, error) {

	category, err := s.repository.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}

	return category, nil
}git sta