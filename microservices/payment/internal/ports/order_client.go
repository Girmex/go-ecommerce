package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/order/proto"
)

type OrderClient interface {
	UpdateOrderStatus(
		ctx context.Context,
		orderID uint,
		status proto.OrderStatus,
	) (*proto.Order, error)
}