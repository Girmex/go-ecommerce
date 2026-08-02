package grpc

import (
	"errors"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
	proto "github.com/Girmex/go-ecommerce/microservices/order/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusError(err error) error {
	switch {

	case errors.Is(err, domain.ErrOrderNotFound):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case errors.Is(err, domain.ErrForbidden):
		return status.Error(
			codes.PermissionDenied,
			err.Error(),
		)
	case errors.Is(err, domain.ErrInsufficientStock):
		return status.Error(
			codes.FailedPrecondition,
			err.Error(),
		)

	default:
		return status.Error(
			codes.Internal,
			err.Error(),
		)
	}
}

func toDomainOrderStatus(
	status proto.OrderStatus,
) (domain.OrderStatus, error) {

	switch status {
	case proto.OrderStatus_ORDER_STATUS_PENDING:
		return domain.OrderPending, nil

	case proto.OrderStatus_ORDER_STATUS_AWAITING_PAYMENT:
		return domain.OrderAwaitingPayment, nil

	case proto.OrderStatus_ORDER_STATUS_PAID:
		return domain.OrderPaid, nil

	case proto.OrderStatus_ORDER_STATUS_CANCELLED:
		return domain.OrderCancelled, nil

	case proto.OrderStatus_ORDER_STATUS_SHIPPED:
		return domain.OrderShipped, nil

	case proto.OrderStatus_ORDER_STATUS_DELIVERED:
		return domain.OrderDelivered, nil

	default:
		return "", errors.New("invalid order status")
	}
}
