package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/port"
)

var _ port.PaymentGateway = (*ChapaGateway)(nil)

type ChapaGateway struct {
	secretKey string
	baseURL   string
	client    *http.Client
}

func NewChapaGateway(secretKey string) *ChapaGateway {
	return &ChapaGateway{
		secretKey: secretKey,
		baseURL:   "https://api.chapa.co/v1",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type chapaInitializeRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number,omitempty"`
	TxRef       string `json:"tx_ref"`
	CallbackURL string `json:"callback_url"`
	ReturnURL   string `json:"return_url"`
}

type chapaInitializeResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

func (g *ChapaGateway) Charge(ctx context.Context, input port.GatewayChargeInput) (port.GatewayChargeResult, error) {
	if input.Amount <= 0 {
		return port.GatewayChargeResult{Approved: false}, errors.New("amount must be greater than zero")
	}

	reqBody := chapaInitializeRequest{
		Amount:      fmt.Sprintf("%.2f", input.Amount),
		Currency:    input.Currency,
		Email:       input.Email,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		TxRef:       input.Reference,
		CallbackURL: input.CallbackURL,
		ReturnURL:   input.ReturnURL,
	}

	if reqBody.Currency == "" {
		reqBody.Currency = "ETB"
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return port.GatewayChargeResult{Approved: false}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/transaction/initialize", bytes.NewBuffer(reqBytes))
	if err != nil {
		return port.GatewayChargeResult{Approved: false}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return port.GatewayChargeResult{Approved: false}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return port.GatewayChargeResult{Approved: false}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return port.GatewayChargeResult{Approved: false}, fmt.Errorf("chapa API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chapaResp chapaInitializeResponse
	if err := json.Unmarshal(body, &chapaResp); err != nil {
		return port.GatewayChargeResult{Approved: false}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chapaResp.Status != "success" {
		return port.GatewayChargeResult{Approved: false}, fmt.Errorf("chapa initialization failed: %s", chapaResp.Message)
	}

	return port.GatewayChargeResult{
		Approved:    true,
		TxnRef:      input.Reference,
		CheckoutURL: chapaResp.Data.CheckoutURL,
	}, nil
}
