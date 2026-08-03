package auth

import (
	authproto "github.com/Girmex/go-ecommerce/microservices/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	authproto.AuthServiceClient
}

func New(addr string) (*Client, error) {

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		AuthServiceClient: authproto.NewAuthServiceClient(conn),
	}, nil
}