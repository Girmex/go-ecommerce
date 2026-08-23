package port

import "context"

type ProductInfo struct {
	ID    string
	Name  string
	Price float64
	Stock int
}

type ProductClient interface {
	GetProduct(ctx context.Context, productID string) (*ProductInfo, error)
	ReserveStock(ctx context.Context, productID string, quantity int) error
}