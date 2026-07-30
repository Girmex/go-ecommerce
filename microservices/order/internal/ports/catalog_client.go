package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/catalog/proto"
)

type CatalogClient interface {
	GetProduct(
		ctx context.Context,
		productID uint,
	) (*proto.Product, error)
	DecreaseProductStock(
		ctx context.Context,
		productID uint,
		quantity uint32,
	) (*proto.Product, error)
}