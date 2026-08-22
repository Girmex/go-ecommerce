package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/port"
)

// ============================================================
// Generic HTTP Client
// ============================================================

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Get(
	ctx context.Context,
	path string,
	result any,
) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+path,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")

	return c.do(req, result)
}

func (c *Client) Post(
	ctx context.Context,
	path string,
	body any,
	result any,
) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+path,
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.do(req, result)
}

func (c *Client) do(
	req *http.Request,
	result any,
) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf(
			"downstream service returned status %d",
			resp.StatusCode,
		)
	}

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// ============================================================
// User Client
// ============================================================

type UserHTTPClient struct {
	c *Client
}

func NewUserHTTPClient(baseURL string) *UserHTTPClient {
	return &UserHTTPClient{
		c: New(baseURL),
	}
}

type userDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *UserHTTPClient) GetUser(
	ctx context.Context,
	userID string,
) (*port.UserInfo, error) {

	var dto userDTO

	if err := c.c.Get(
		ctx,
		"/users/"+userID,
		&dto,
	); err != nil {
		return nil, fmt.Errorf("user-ms: %w", err)
	}

	return &port.UserInfo{
		ID:    dto.ID,
		Name:  dto.Name,
		Email: dto.Email,
	}, nil
}

// ============================================================
// Product Client
// ============================================================

type ProductHTTPClient struct {
	c *Client
}

func NewProductHTTPClient(baseURL string) *ProductHTTPClient {
	return &ProductHTTPClient{
		c: New(baseURL),
	}
}

type productDTO struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func (c *ProductHTTPClient) GetProduct(
	ctx context.Context,
	productID string,
) (*port.ProductInfo, error) {

	var dto productDTO

	if err := c.c.Get(
		ctx,
		"/products/"+productID,
		&dto,
	); err != nil {
		return nil, fmt.Errorf("product-ms: %w", err)
	}

	return &port.ProductInfo{
		ID:    dto.ID,
		Name:  dto.Name,
		Price: dto.Price,
		Stock: dto.Stock,
	}, nil
}

func (c *ProductHTTPClient) ReserveStock(
	ctx context.Context,
	productID string,
	quantity int,
) error {

	body := map[string]int{
		"quantity": quantity,
	}

	if err := c.c.Post(
		ctx,
		"/products/"+productID+"/reserve",
		body,
		nil,
	); err != nil {
		return fmt.Errorf("product-ms: %w", err)
	}

	return nil
}

// ============================================================
// Payment Client
// ============================================================

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
	userID string,
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