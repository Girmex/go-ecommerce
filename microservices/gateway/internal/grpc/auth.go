package grpc

import (
	authproto "github.com/Girmex/go-ecommerce/microservices/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth authproto.AuthServiceClient
}

func NewClients(authAddr string) (*Clients, error) {

	conn, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Auth: authproto.NewAuthServiceClient(conn),
	}, nil
}