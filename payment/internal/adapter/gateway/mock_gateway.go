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
	input port.GatewayChargeInput,
) (port.GatewayChargeResult, error) {

	if input.Amount <= 0 {
		return port.GatewayChargeResult{}, nil
	}

	return port.GatewayChargeResult{
		TxnRef:      "txn_" + uuid.NewString(),
		CheckoutURL: "",
	}, nil
}