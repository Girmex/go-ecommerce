package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/port"
)

type orderService struct {
	repo          port.OrderRepository
	productClient port.ProductClient
	paymentClient port.PaymentClient
}

func NewOrderService(
	repo port.OrderRepository,
	productClient port.ProductClient,
	paymentClient port.PaymentClient,
) port.OrderService {
	return &orderService{
		repo:          repo,
		productClient: productClient,
		paymentClient: paymentClient,
	}
}

func (s *orderService) PlaceOrder(ctx context.Context, in port.CreateOrderInput) (*domain.Order, error) {
	if len(in.Items) == 0 {
		return nil, domain.ErrEmptyOrder
	}

	//  Validate each product, snapshot its price, and reserve stock
	//    (call out to product-ms).
	items := make([]domain.OrderItem, 0, len(in.Items))
	for _, reqItem := range in.Items {
		product, err := s.productClient.GetProduct(ctx, reqItem.ProductID)
		if err != nil {
			return nil, domain.ErrProductNotFound
		}
		if err := s.productClient.ReserveStock(ctx, reqItem.ProductID, reqItem.Quantity); err != nil {
			return nil, err
		}
		items = append(items, domain.OrderItem{
			ProductID: product.ID,
			Quantity:  reqItem.Quantity,
			UnitPrice: product.Price,
		})
	}

	now := time.Now().UTC()
	order := &domain.Order{
		ID:        uuid.NewString(),
		UserID:    in.UserID,
		Items:     items,
		Status:    domain.StatusCreated,
		CreatedAt: now,
		UpdatedAt: now,
	}
	order.ComputeTotal()

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	fmt.Println("ABOUT TO CREATE ORDER:", order.ID)

	// 3. Charge payment (call out to payment-ms). A decline doesn't
	//    fail the whole request — the order is kept with a
	//    payment_failed status so the client can retry payment.
	paymentResult, err := s.paymentClient.Charge(ctx, order.ID, order.UserID, order.Total, in.PaymentMethod)
	order.UpdatedAt = time.Now().UTC()
	if err != nil || paymentResult == nil || paymentResult.Status != "paid" {
		order.Status = domain.StatusPaymentFailed
		_ = s.repo.Update(ctx, order)
		return order, nil
	}

	order.Status = domain.StatusPaid
	order.PaymentID = paymentResult.ID
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) Get(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) ListByUser(ctx context.Context, userID uint) ([]*domain.Order, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *orderService) Cancel(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.Status = domain.StatusCancelled
	order.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}
