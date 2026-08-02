package domain

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidOrderStatus = errors.New("invalid order status")
)