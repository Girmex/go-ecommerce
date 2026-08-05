package payment

import (

	paymentproto "github.com/Girmex/go-ecommerce/microservices/payment/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	paymentproto.PaymentServiceClient
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
		PaymentServiceClient: paymentproto.NewPaymentServiceClient(conn),
	}, nil
}