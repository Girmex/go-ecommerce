package client

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/port"
)

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

	fmt.Printf(
		"calling Product-MS: %s/products/%s\n",
		c.c.baseURL,
		productID,
	)

	var dto productDTO

	if err := c.c.Get(
		ctx,
		"/products/"+productID,
		&dto,
	); err != nil {
		return nil, fmt.Errorf(
			"product-ms GetProduct(%s): %w",
			productID,
			err,
		)
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
