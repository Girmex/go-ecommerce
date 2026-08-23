package http

import "github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/domain"

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid" example:"b3c1..."`
	Quantity  int    `json:"quantity" validate:"required,gt=0" example:"2"`
}

type CreateOrderRequest struct {
	Items         []CreateOrderItemRequest  `json:"items" validate:"required,min=1,dive"`
	PaymentMethod string                    `json:"payment_method" validate:"required,oneof=card wallet bank_transfer" example:"card"`
}

type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID        string              `json:"id"`
	UserID    uint                `json:"user_id"`
	Items     []OrderItemResponse `json:"items"`
	Total     float64             `json:"total"`
	Status    string              `json:"status"`
	PaymentID string              `json:"payment_id,omitempty"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

func toOrderResponse(o *domain.Order) OrderResponse {
	items := make([]OrderItemResponse, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, OrderItemResponse{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			Subtotal:  it.Subtotal(),
		})
	}
	return OrderResponse{
		ID:        o.ID,
		UserID:    o.UserID,
		Items:     items,
		Total:     o.Total,
		Status:    string(o.Status),
		PaymentID: o.PaymentID,
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
