package port

import "context"

type GatewayChargeResult struct {
	Approved bool
	TxnRef   string
}

type PaymentGateway interface {
	Charge(
		ctx context.Context,
		amount float64,
		method string,
	) (GatewayChargeResult, error)
}