package port

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) error
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
	List(ctx context.Context) ([]*domain.Payment, error)
	Update(ctx context.Context, p *domain.Payment) error
}