package domain

import (
	"errors"
	"time"
)

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrEmptyOrder     = errors.New("order must contain at least one item")
	ErrUserNotFound   = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")
)

type OrderStatus string

const (
	StatusCreated       OrderStatus = "created"
	StatusPaymentFailed OrderStatus = "payment_failed"
	StatusPaid          OrderStatus = "paid"
	StatusCancelled     OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func (i OrderItem) Subtotal() float64 {
	return i.UnitPrice * float64(i.Quantity)
}

type Order struct {
	ID        string      `json:"id"`
	UserID    uint        `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    OrderStatus `json:"status"`
	PaymentID string      `json:"payment_id,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (o *Order) ComputeTotal() {
	var total float64
	for _, item := range o.Items {
		total += item.Subtotal()
	}
	o.Total = total
}
