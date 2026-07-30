package domain

import "errors"

var (
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrPaymentAlreadyCompleted = errors.New("payment already completed")
	ErrPaymentAlreadyFailed    = errors.New("payment already failed")
)
