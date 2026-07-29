package application

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/ports"
)

type OrderService struct {
	repository ports.OrderRepository
}

func NewOrderService(
	orderRepository ports.OrderRepository,
) *OrderService {
	return &OrderService{
		repository: orderRepository,
	}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	order *domain.Order,
) error {

	return s.repository.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	orderID uint,
) (*domain.Order, error) {

	return s.repository.GetOrder(ctx, orderID)
}

func (s *OrderService) GetOrdersByUser(
	ctx context.Context,
	userID uint,
) ([]domain.Order, error) {

	return s.repository.GetOrdersByUser(ctx, userID)
}

func (s *OrderService) CancelOrder(
	ctx context.Context,
	orderID uint,
) error {

	order, err := s.repository.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = domain.OrderCancelled

	return s.repository.UpdateOrder(ctx, order)
}