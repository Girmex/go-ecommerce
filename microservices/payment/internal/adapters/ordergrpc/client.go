package ordergrpc

import (
	"context"

	orderproto "github.com/Girmex/go-ecommerce/microservices/order/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client orderproto.OrderServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		client: orderproto.NewOrderServiceClient(conn),
	}
}

func (c *Client) UpdateOrderStatus(
	ctx context.Context,
	orderID uint,
	status orderproto.OrderStatus,
) (*orderproto.Order, error) {

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

	return c.client.UpdateOrderStatus(
		ctx,
		&orderproto.UpdateOrderStatusRequest{
			Id:     uint32(orderID),
			Status: status,
		},
	)
}
