package application

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/payment/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/payment/internal/ports"
)

type PaymentService struct {
	repository  ports.PaymentRepository
	orderClient ports.OrderClient
}

func NewPaymentService(
	repository ports.PaymentRepository,
	orderClient ports.OrderClient,
) *PaymentService {
	return &PaymentService{
		repository:  repository,
		orderClient: orderClient,
	}
}

func (s *PaymentService) CreatePayment(
	ctx context.Context,
	orderID uint,
	userID uint,
	amount float64,
) (*domain.Payment, error) {

	payment := &domain.Payment{
		OrderID: orderID,
		UserID:  userID,
		Amount:  amount,
		Status:  domain.PaymentPending,
	}

	err := s.repository.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) GetPayment(
	ctx context.Context,
	id uint,
) (*domain.Payment, error) {

	return s.repository.GetPayment(ctx, id)
}

func (s *PaymentService) CompletePayment(
	ctx context.Context,
	id uint,
) (*domain.Payment, error) {

	payment, err := s.repository.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.orderClient.MarkOrderAsPaid(
		ctx,
		payment.OrderID,
	); err != nil {
		return nil, err
	}

	if err := payment.Complete(); err != nil {
		return nil, err
	}

	if err := s.repository.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) FailPayment(
	ctx context.Context,
	id uint,
) (*domain.Payment, error) {

	payment, err := s.repository.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.orderClient.CancelOrder(
		ctx,
		payment.OrderID,
	); err != nil {
		return nil, err
	}

	if err := payment.Fail(); err != nil {
		return nil, err
	}

	if err := s.repository.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
