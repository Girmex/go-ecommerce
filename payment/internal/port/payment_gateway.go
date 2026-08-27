package port

import "context"

type GatewayChargeInput struct {
	Amount      float64
	Currency    string
	Reference   string
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	CallbackURL string
	ReturnURL   string
}

type GatewayChargeResult struct {
	Approved    bool
	TxnRef      string
	CheckoutURL string
}

type PaymentGateway interface {
	Charge(
		ctx context.Context,
		input GatewayChargeInput,
	) (GatewayChargeResult, error)
}