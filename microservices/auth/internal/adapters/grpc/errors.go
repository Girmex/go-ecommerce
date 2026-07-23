package grpc

import (
	"errors"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusError(err error) error {

	switch {

	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(
			codes.AlreadyExists,
			err.Error(),
		)

	case errors.Is(err, domain.	ErrInvalidCredentials):
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