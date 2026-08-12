package authgrpc

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce/microservices/auth/proto"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
	"google.golang.org/grpc"
)

var _ ports.UserClient = (*Client)(nil)

type Client struct {
	client proto.AuthServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		client: proto.NewAuthServiceClient(conn),
	}
}

func (c *Client) GetUser(
	ctx context.Context,
	userID uint,
) (*ports.User, error) {
	resp, err := c.client.GetUser(
		ctx,
		&proto.GetUserRequest{
			Id: uint32(userID),
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get user from auth service: %w",
			err,
		)
	}

	return &ports.User{
		ID:    uint(resp.Id),
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}
