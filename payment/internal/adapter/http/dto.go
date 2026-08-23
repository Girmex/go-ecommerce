package http

import (
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/domain"
)

type ChargeRequest struct {
	OrderID string  `json:"order_id" validate:"required,uuid" example:"b3c1d7e2-1234-4567-8901-abcdef123456"`
	UserID  uint    `json:"user_id" validate:"required" example:"1234567890"`
	Amount  float64 `json:"amount" validate:"required,gt=0" example:"59.98"`
	Method  string  `json:"method" validate:"required,oneof=card wallet bank_transfer" example:"card"`
}

type PaymentResponse struct {
	ID            string  `json:"id"`
	OrderID       string  `json:"order_id"`
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	Method        string  `json:"method"`
	Status        string  `json:"status"`
	GatewayTxnRef string  `json:"gateway_txn_ref,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func toPaymentResponse(
	p *domain.Payment,
) PaymentResponse {
	return PaymentResponse{
		ID:            p.ID,
		OrderID:       p.OrderID,
		UserID:        p.UserID,
		Amount:        p.Amount,
		Method:        p.Method,
		Status:        string(p.Status),
		GatewayTxnRef: p.GatewayTxnRef,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
