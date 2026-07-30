package grpc

import (
	"errors"

	"github.com/Girmex/go-ecommerce/microservices/payment/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusError(err error) error {

	switch {

	case errors.Is(err, domain.ErrPaymentNotFound):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case errors.Is(err, domain.ErrPaymentAlreadyCompleted):
		return status.Error(
			codes.FailedPrecondition,
			err.Error(),
		)

	case errors.Is(err, domain.ErrPaymentAlreadyFailed):
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