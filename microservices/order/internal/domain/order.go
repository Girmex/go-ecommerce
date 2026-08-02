package domain

import (
	"time"
)

type Order struct {
	ID         uint
	UserID     uint
	Status     OrderStatus
	TotalPrice float64
	Items      []OrderItem
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (o *Order) MarkAsPaid() error {

	if o.Status != OrderPending {
		return ErrInvalidOrderStatus
	}

	o.Status = OrderPaid

	return nil
}

func (o *Order) Cancel() error {

	if o.Status == OrderDelivered {
		return ErrInvalidOrderStatus
	}

	if o.Status == OrderCancelled {
		return ErrInvalidOrderStatus
	}

	o.Status = OrderCancelled

	return nil
}
