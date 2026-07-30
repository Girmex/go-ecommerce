package domain

import "time"

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID uint

	OrderID uint
	UserID uint

	Amount float64

	Status PaymentStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}