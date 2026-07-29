package grpc

import (
	proto "github.com/Girmex/go-ecommerce/microservices/order/proto"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
)

func toProtoStatus(status domain.OrderStatus) proto.OrderStatus {
	switch status {
	case domain.OrderPending:
		return proto.OrderStatus_ORDER_STATUS_PENDING

	case domain.OrderAwaitingPayment:
		return proto.OrderStatus_ORDER_STATUS_AWAITING_PAYMENT

	case domain.OrderPaid:
		return proto.OrderStatus_ORDER_STATUS_PAID

	case domain.OrderCancelled:
		return proto.OrderStatus_ORDER_STATUS_CANCELLED

	case domain.OrderShipped:
		return proto.OrderStatus_ORDER_STATUS_SHIPPED

	case domain.OrderDelivered:
		return proto.OrderStatus_ORDER_STATUS_DELIVERED

	default:
		return proto.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func toProtoOrder(order *domain.Order) *proto.Order {

	items := make([]*proto.OrderItem, len(order.Items))

	for i, item := range order.Items {
		items[i] = &proto.OrderItem{
			Id:        uint32(item.ID),
			ProductId: uint32(item.ProductID),
			Quantity:  uint32(item.Quantity),
			UnitPrice: item.UnitPrice,
		}
	}

	return &proto.Order{
		Id:         uint32(order.ID),
		UserId:     uint32(order.UserID),
		Status:     toProtoStatus(order.Status),
		Items:      items,
		TotalPrice: order.TotalPrice,
	}
}

func toDomainOrderItems(
	items []*proto.OrderItemRequest,
) []domain.OrderItem {

	result := make([]domain.OrderItem, len(items))

	for i, item := range items {
		result[i] = domain.OrderItem{
			ProductID: uint(item.ProductId),
			Quantity:  uint(item.Quantity),
		}
	}

	return result
}
