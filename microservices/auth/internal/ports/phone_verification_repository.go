package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
)

type PhoneVerificationRepository interface {
	Create(
		ctx context.Context,
		verification *domain.PhoneVerification,
	) error

	GetLatestByUserID(
		ctx context.Context,
		userID uint,
	) (*domain.PhoneVerification, error)

	MarkUsed(
		ctx context.Context,
		verification *domain.PhoneVerification,
	) error
}