package ports

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/payment/internal/domain"
)
type PaymentRepository interface {

	CreatePayment(
		ctx context.Context,
		payment *domain.Payment,
	) error

	GetPayment(
		ctx context.Context,
		id uint,
	) (*domain.Payment,error)

	UpdatePayment(
		ctx context.Context,
		payment *domain.Payment,
	) error
}