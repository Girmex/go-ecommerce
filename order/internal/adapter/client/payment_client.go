package client

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/port"
)

type PaymentHTTPClient struct {
	c *Client
}

func NewPaymentHTTPClient(baseURL string) *PaymentHTTPClient {
	return &PaymentHTTPClient{
		c: New(baseURL),
	}
}

type paymentDTO struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (c *PaymentHTTPClient) Charge(
	ctx context.Context,
	orderID string,
	userID uint,
	amount float64,
	method string,
) (*port.PaymentResult, error) {

	body := map[string]any{
		"order_id": orderID,
		"user_id":  userID,
		"amount":   amount,
		"method":   method,
	}

	var dto paymentDTO

	if err := c.c.Post(
		ctx,
		"/payments",
		body,
		&dto,
	); err != nil {
		return nil, fmt.Errorf("payment-ms: %w", err)
	}

	return &port.PaymentResult{
		ID:     dto.ID,
		Status: dto.Status,
	}, nil
}