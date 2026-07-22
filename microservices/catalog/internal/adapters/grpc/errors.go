package grpc

import (
	"errors"

	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusError(err error) error {

	switch {

	case errors.Is(err, domain.ErrCategoryNotFound):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case errors.Is(err, domain.ErrProductNotFound):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(
			codes.PermissionDenied,
			err.Error(),
		)

	default:
		return status.Error(
			codes.Internal,
			err.Error(),
		)
	}
}