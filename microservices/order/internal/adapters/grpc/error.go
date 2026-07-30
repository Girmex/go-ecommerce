package grpc

import (
	"errors"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"

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
