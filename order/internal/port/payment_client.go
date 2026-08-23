package port

import "context"

type PaymentResult struct {
	ID     string
	Status string
}

type PaymentClient interface {
	Charge(
		ctx context.Context,
		orderID string,
		userID string,
		amount float64,
		method string,
	) (*PaymentResult, error)
}