package authgrpc

import (
	"context"

	authproto "github.com/Girmex/go-ecommerce/microservices/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client authproto.AuthServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		client: authproto.NewAuthServiceClient(conn),
	}
}

func (c *Client) GetUser(
	ctx context.Context,
	userID uint,
) (*authproto.User, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		authHeaders := md.Get("authorization")
		if len(authHeaders) > 0 {
			ctx = metadata.AppendToOutgoingContext(
				ctx,
				"authorization",
				authHeaders[0],
			)
		}
	}

	return c.client.GetUser(
		ctx,
		&authproto.GetUserRequest{
			Id: uint32(userID),
		},
	)
}
