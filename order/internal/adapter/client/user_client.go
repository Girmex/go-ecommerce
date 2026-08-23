package client

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/port"
)

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