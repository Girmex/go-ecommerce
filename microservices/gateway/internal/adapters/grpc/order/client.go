package order

import (
	orderproto "github.com/Girmex/go-ecommerce/microservices/order/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	orderproto.OrderServiceClient
}

func New(addr string) (*Client, error) {

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		OrderServiceClient: orderproto.NewOrderServiceClient(conn),
	}, nil
}