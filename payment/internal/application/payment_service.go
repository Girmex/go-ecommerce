package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/port"
)

type PaymentService struct {
	repository port.PaymentRepository
	gateway    port.PaymentGateway
}

func NewPaymentService(
	repository port.PaymentRepository,
	gateway port.PaymentGateway,
) port.PaymentService {
	return &PaymentService{
		repository: repository,
		gateway:    gateway,
	}
}

func (s *PaymentService) Charge(
	ctx context.Context,
	input port.ChargeInput,
) (*port.ChargeResult, error) {

	now := time.Now().UTC()

	payment := &domain.Payment{
		ID:        uuid.NewString(),
		OrderID:   input.OrderID,
		UserID:    input.UserID,
		Amount:    input.Amount,
		Method:    input.Method,
		Status:    domain.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the payment as PENDING before contacting the gateway.
	if err := s.repository.Create(ctx, payment); err != nil {
		return nil, err
	}

	result, err := s.gateway.Charge(
		ctx,
		port.GatewayChargeInput{
			Amount:      input.Amount,
			Currency:    input.Currency,
			Reference:   payment.ID,
			Email:       input.Email,
			FirstName:   input.FirstName,
			LastName:    input.LastName,
			PhoneNumber: input.PhoneNumber,
			CallbackURL: input.CallbackURL,
			ReturnURL:   input.ReturnURL,
		},
	)
	if err != nil {
		return nil, err
	}

	// Gateway initialization does NOT mean the payment was completed.
	// The payment remains PENDING until verification/webhook processing.
	payment.GatewayTxnRef = result.TxnRef
	payment.UpdatedAt = time.Now().UTC()

	if err := s.repository.Update(ctx, payment); err != nil {
		return nil, err
	}

	return &port.ChargeResult{
		Payment:     payment,
		CheckoutURL: result.CheckoutURL,
	}, nil
}

func (s *PaymentService) Get(
	ctx context.Context,
	id string,
) (*domain.Payment, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *PaymentService) List(
	ctx context.Context,
) ([]*domain.Payment, error) {
	return s.repository.List(ctx)
}

func (s *PaymentService) Refund(
	ctx context.Context,
	id string,
) (*domain.Payment, error) {

	payment, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	payment.Status = domain.StatusRefunded
	payment.UpdatedAt = time.Now().UTC()

	if err := s.repository.Update(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}