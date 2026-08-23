package gateway

import (
	"context"

	"github.com/google/uuid"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/port"
)

var _ port.PaymentGateway = (*MockGateway)(nil)

type MockGateway struct{}

func NewMockGateway() *MockGateway {
	return &MockGateway{}
}

func (g *MockGateway) Charge(
	_ context.Context,
	amount float64,
	_ string,
) (port.GatewayChargeResult, error) {

	if amount <= 0 {
		return port.GatewayChargeResult{
			Approved: false,
		}, nil
	}

	return port.GatewayChargeResult{
		Approved: true,
		TxnRef:   "txn_" + uuid.NewString(),
	}, nil
}