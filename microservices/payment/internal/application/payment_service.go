package application

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/payment/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/payment/internal/ports"
	"github.com/Girmex/go-ecommerce/microservices/order/proto"
	
)

type PaymentService struct {
	repository ports.PaymentRepository
	order      ports.OrderClient
}

func NewPaymentService(
	repository ports.PaymentRepository,
	orderClient ports.OrderClient,
) *PaymentService {
	return &PaymentService{
		repository: repository,
		order:      orderClient,
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

	if payment.Status == domain.PaymentSuccess {
		return nil, domain.ErrPaymentAlreadyCompleted
	}

	if payment.Status == domain.PaymentFailed {
		return nil, domain.ErrPaymentAlreadyFailed
	}

	payment.Status = domain.PaymentSuccess

	if err := s.repository.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	_, err = s.order.UpdateOrderStatus(
		ctx,
		payment.OrderID,
		proto.OrderStatus_ORDER_STATUS_PAID,
	)
	if err != nil {
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

	if payment.Status == domain.PaymentSuccess {
		return nil, domain.ErrPaymentAlreadyCompleted
	}

	if payment.Status == domain.PaymentFailed {
		return nil, domain.ErrPaymentAlreadyFailed
	}

	payment.Status = domain.PaymentFailed

	if err := s.repository.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
