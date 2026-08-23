package domain

import (
	"errors"
	"time"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidPayment  = errors.New("invalid payment")
	ErrPaymentDeclined = errors.New("payment declined")
)

type PaymentStatus string

const (
	StatusPending  PaymentStatus = "pending"
	StatusPaid     PaymentStatus = "paid"
	StatusDeclined PaymentStatus = "declined"
	StatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID             string        `json:"id"`
	OrderID        string        `json:"order_id"`
	UserID         string        `json:"user_id"`
	Amount         float64       `json:"amount"`
	Method         string        `json:"method"`
	Status         PaymentStatus `json:"status"`
	GatewayTxnRef  string        `json:"gateway_txn_ref,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}