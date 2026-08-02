package application

import (
	"context"
	"log"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/ports"
)

type OrderService struct {
	repository ports.OrderRepository
	catalog    ports.CatalogClient
}

func NewOrderService(
	orderRepository ports.OrderRepository,
	catalogClient ports.CatalogClient,
) *OrderService {
	return &OrderService{
		repository: orderRepository,
		catalog:    catalogClient,
	}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	input dto.CreateOrderInput,
	userID uint32,
) (*domain.Order, error) {

	items := make([]domain.OrderItem, 0, len(input.Items))

	var totalPrice float64

	for _, item := range input.Items {

		product, err := s.catalog.GetProduct(
			ctx,
			uint(item.ProductID),
		)
		if err != nil {
			return nil, err
		}

		// Validate stock
		if product.Stock < item.Quantity {
			return nil, domain.ErrInsufficientStock
		}

		// Decrease stock in Catalog Service
		_, err = s.catalog.DecreaseProductStock(
			ctx,
			uint(product.Id),
			item.Quantity,
		)
		if err != nil {
			return nil, err
		}

		orderItem := domain.OrderItem{
			ProductID: uint(product.Id),
			Quantity:  uint(item.Quantity),
			UnitPrice: product.Price,
		}

		items = append(items, orderItem)

		totalPrice += product.Price * float64(item.Quantity)
	}

	order := &domain.Order{
		UserID:     uint(userID),
		Status:     domain.OrderPending,
		Items:      items,
		TotalPrice: totalPrice,
	}

	err := s.repository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	log.Printf("TOTAL: %.2f", order.TotalPrice)

	for _, item := range order.Items {
		log.Printf(
			"ITEM -> product=%d unit_price=%.2f quantity=%d",
			item.ProductID,
			item.UnitPrice,
			item.Quantity,
		)
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

func (s *OrderService) UpdateOrderStatus(
	ctx context.Context,
	id uint,
	status domain.OrderStatus,
) (*domain.Order, error) {

	return s.repository.UpdateOrderStatus(
		ctx,
		id,
		status,
	)
}
