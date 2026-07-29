package application

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/dto"
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
	input dto.CreateOrderInput,
	userID uint32,
) (*domain.Order, error) {

	items := make([]domain.OrderItem, 0, len(input.Items))

	for _, item := range input.Items {
		items = append(items, domain.OrderItem{
			ProductID: uint(item.ProductID),
			Quantity:  uint(item.Quantity),
		})
	}

	order := &domain.Order{
		UserID: uint(userID),
		Status: domain.OrderPending,
		Items:  items,
	}

	err := s.repository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	orderID uint,
) (*domain.Order, error) {

	return s.repository.GetOrder(ctx, orderID)
}

func (s *OrderService) GetOrdersByUser(
	ctx context.Context,
	userID uint32,
) ([]domain.Order, error) {

	return s.repository.GetOrdersByUser(ctx, uint(userID))
}

func (s *OrderService) CancelOrder(
	ctx context.Context,
	orderID uint,
	userID uint32,
) error {

	order, err := s.repository.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != uint(userID) {
		return domain.ErrForbidden
	}

	order.Status = domain.OrderCancelled

	return s.repository.UpdateOrder(ctx, order)
}