package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/domain"
	port "github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/port"
)

// Mock ProductRepository

type mockProductRepository struct {
	products map[string]*domain.Product
	err      error
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (m *mockProductRepository) Create(
	ctx context.Context,
	p *domain.Product,
) error {
	if m.err != nil {
		return m.err
	}

	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.Product, error) {
	if m.err != nil {
		return nil, m.err
	}

	p, ok := m.products[id]
	if !ok {
		return nil, domain.ErrProductNotFound
	}

	return p, nil
}

func (m *mockProductRepository) List(
	ctx context.Context,
) ([]*domain.Product, error) {
	if m.err != nil {
		return nil, m.err
	}

	products := make([]*domain.Product, 0, len(m.products))

	for _, p := range m.products {
		products = append(products, p)
	}

	return products, nil
}

func (m *mockProductRepository) Update(
	ctx context.Context,
	p *domain.Product,
	userID uint,
) error {
	if m.err != nil {
		return m.err
	}

	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepository) Delete(
	ctx context.Context,
	id string,
	userID uint,
) error {
	if m.err != nil {
		return m.err
	}

	if _, ok := m.products[id]; !ok {
		return domain.ErrProductNotFound
	}

	delete(m.products, id)

	return nil
}

var _ port.ProductRepository = (*mockProductRepository)(nil)

// Helper

func testProduct(userID uint) *domain.Product {
	now := time.Now().UTC()

	return &domain.Product{
		ID:          "product-1",
		UserID:      userID,
		Name:        "Test Product",
		Description: "Test description",
		Price:       100,
		Stock:       10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// --------------------------------------------------
// Create
// --------------------------------------------------

func TestCreate(t *testing.T) {
	repository := newMockProductRepository()
	service := NewProductService(repository)

	product, err := service.Create(
		context.Background(),
		port.CreateProductInput{
			UserID:      1,
			Name:        "Test Product",
			Description: "Test description",
			Price:       100,
			Stock:       10,
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}

	if product.ID == "" {
		t.Error("expected product ID to be generated")
	}

	if product.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", product.UserID)
	}

	if product.Name != "Test Product" {
		t.Errorf("expected name Test Product, got %s", product.Name)
	}

	if product.Description != "Test description" {
		t.Errorf(
			"expected description Test description, got %s",
			product.Description,
		)
	}

	if product.Price != 100 {
		t.Errorf("expected price 100, got %v", product.Price)
	}

	if product.Stock != 10 {
		t.Errorf("expected stock 10, got %d", product.Stock)
	}

	if product.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if product.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	if _, ok := repository.products[product.ID]; !ok {
		t.Error("expected product to be stored in repository")
	}
}

func TestCreate_RepositoryError(t *testing.T) {
	repository := newMockProductRepository()
	repository.err = errors.New("repository error")

	service := NewProductService(repository)

	_, err := service.Create(
		context.Background(),
		port.CreateProductInput{
			UserID:      1,
			Name:        "Test Product",
			Description: "Test description",
			Price:       100,
			Stock:       10,
		},
	)

	if !errors.Is(err, repository.err) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// Get
// --------------------------------------------------

func TestGet(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	repository.products[product.ID] = product

	service := NewProductService(repository)

	result, err := service.Get(
		context.Background(),
		product.ID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected product, got nil")
	}

	if result.ID != product.ID {
		t.Errorf(
			"expected ID %s, got %s",
			product.ID,
			result.ID,
		)
	}

	if result.Name != "Test Product" {
		t.Errorf(
			"expected name Test Product, got %s",
			result.Name,
		)
	}
}

func TestGet_NotFound(t *testing.T) {
	repository := newMockProductRepository()
	service := NewProductService(repository)

	_, err := service.Get(
		context.Background(),
		"unknown-product",
	)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Fatalf(
			"expected ErrProductNotFound, got %v",
			err,
		)
	}
}

func TestGet_RepositoryError(t *testing.T) {
	repository := newMockProductRepository()
	repository.err = errors.New("repository error")

	service := NewProductService(repository)

	_, err := service.Get(
		context.Background(),
		"product-1",
	)

	if !errors.Is(err, repository.err) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// List
// --------------------------------------------------

func TestList(t *testing.T) {
	repository := newMockProductRepository()

	product1 := testProduct(1)

	product2 := &domain.Product{
		ID:          "product-2",
		UserID:      2,
		Name:        "Second Product",
		Description: "Second description",
		Price:       200,
		Stock:       20,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	repository.products[product1.ID] = product1
	repository.products[product2.ID] = product2

	service := NewProductService(repository)

	products, err := service.List(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 2 {
		t.Fatalf(
			"expected 2 products, got %d",
			len(products),
		)
	}
}

func TestList_Empty(t *testing.T) {
	repository := newMockProductRepository()
	service := NewProductService(repository)

	products, err := service.List(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 0 {
		t.Fatalf(
			"expected 0 products, got %d",
			len(products),
		)
	}
}

func TestList_RepositoryError(t *testing.T) {
	repository := newMockProductRepository()
	repository.err = errors.New("repository error")

	service := NewProductService(repository)

	_, err := service.List(context.Background())

	if !errors.Is(err, repository.err) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// Update
// --------------------------------------------------

func TestUpdate(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	repository.products[product.ID] = product

	service := NewProductService(repository)

	newStock := 20

	result, err := service.Update(
		context.Background(),
		product.ID,
		port.UpdateProductInput{
			UserID:      1,
			Name:        "Updated Product",
			Description: "Updated description",
			Price:       150,
			Stock:       &newStock,
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected product, got nil")
	}

	if result.Name != "Updated Product" {
		t.Errorf(
			"expected name Updated Product, got %s",
			result.Name,
		)
	}

	if result.Description != "Updated description" {
		t.Errorf(
			"expected description Updated description, got %s",
			result.Description,
		)
	}

	if result.Price != 150 {
		t.Errorf(
			"expected price 150, got %v",
			result.Price,
		)
	}

	if result.Stock != 20 {
		t.Errorf(
			"expected stock 20, got %d",
			result.Stock,
		)
	}

	if result.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repository := newMockProductRepository()
	service := NewProductService(repository)

	_, err := service.Update(
		context.Background(),
		"unknown-product",
		port.UpdateProductInput{
			UserID: 1,
			Name:   "Updated Product",
		},
	)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Fatalf(
			"expected ErrProductNotFound, got %v",
			err,
		)
	}
}

func TestUpdate_NotOwner(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	repository.products[product.ID] = product

	service := NewProductService(repository)

	_, err := service.Update(
		context.Background(),
		product.ID,
		port.UpdateProductInput{
			UserID: 999,
			Name:   "Hacked Product",
		},
	)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Fatalf(
			"expected ErrProductNotFound, got %v",
			err,
		)
	}

	if repository.products[product.ID].Name != "Test Product" {
		t.Error("product should not have been modified")
	}
}

func TestUpdate_RepositoryError(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	repository.products[product.ID] = product

	repository.err = errors.New("repository error")

	service := NewProductService(repository)

	_, err := service.Update(
		context.Background(),
		product.ID,
		port.UpdateProductInput{
			UserID: 1,
			Name:   "Updated Product",
		},
	)

	if !errors.Is(err, repository.err) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// Delete
// --------------------------------------------------

func TestDelete(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	repository.products[product.ID] = product

	service := NewProductService(repository)

	err := service.Delete(
		context.Background(),
		product.ID,
		1,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := repository.products[product.ID]; ok {
		t.Error("expected product to be deleted")
	}
}

func TestDelete_NotFound(t *testing.T) {
	repository := newMockProductRepository()
	service := NewProductService(repository)

	err := service.Delete(
		context.Background(),
		"unknown-product",
		1,
	)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Fatalf(
			"expected ErrProductNotFound, got %v",
			err,
		)
	}
}

func TestDelete_NotOwner(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	repository.products[product.ID] = product

	service := NewProductService(repository)

	err := service.Delete(
		context.Background(),
		product.ID,
		999,
	)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Fatalf(
			"expected ErrProductNotFound, got %v",
			err,
		)
	}

	if _, ok := repository.products[product.ID]; !ok {
		t.Error("product should not have been deleted")
	}
}

func TestDelete_RepositoryError(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	repository.products[product.ID] = product

	repository.err = errors.New("repository error")

	service := NewProductService(repository)

	err := service.Delete(
		context.Background(),
		product.ID,
		1,
	)

	if !errors.Is(err, repository.err) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// ReserveStock
// --------------------------------------------------

func TestReserveStock(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	product.Stock = 10

	repository.products[product.ID] = product

	service := NewProductService(repository)

	result, err := service.ReserveStock(
		context.Background(),
		product.ID,
		3,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected product, got nil")
	}

	if result.Stock != 7 {
		t.Errorf(
			"expected stock 7, got %d",
			result.Stock,
		)
	}
}

func TestReserveStock_NotFound(t *testing.T) {
	repository := newMockProductRepository()
	service := NewProductService(repository)

	_, err := service.ReserveStock(
		context.Background(),
		"unknown-product",
		3,
	)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Fatalf(
			"expected ErrProductNotFound, got %v",
			err,
		)
	}
}

func TestReserveStock_InsufficientStock(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	product.Stock = 2

	repository.products[product.ID] = product

	service := NewProductService(repository)

	_, err := service.ReserveStock(
		context.Background(),
		product.ID,
		5,
	)

	if !errors.Is(err, domain.ErrInsufficientStock) {
		t.Fatalf(
			"expected ErrInsufficientStock, got %v",
			err,
		)
	}
}

func TestReserveStock_RepositoryError(t *testing.T) {
	repository := newMockProductRepository()

	product := testProduct(1)
	product.Stock = 10

	repository.products[product.ID] = product

	repository.err = errors.New("repository error")

	service := NewProductService(repository)

	_, err := service.ReserveStock(
		context.Background(),
		product.ID,
		3,
	)

	if !errors.Is(err, repository.err) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}
