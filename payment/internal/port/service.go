package port

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/domain"
)

type ChargeInput struct {
	OrderID string
	UserID  uint
	Amount  float64
	Method  string
}

type PaymentService interface {
	Charge(ctx context.Context, input ChargeInput) (*domain.Payment, error)
	Get(ctx context.Context, id string) (*domain.Payment, error)
	List(ctx context.Context) ([]*domain.Payment, error)
	Refund(ctx context.Context, id string) (*domain.Payment, error)
}