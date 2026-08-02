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

func (c *Client) MarkOrderAsPaid(
	ctx context.Context,
	orderID uint,
) error {

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

	_, err := c.client.UpdateOrderStatus(
		ctx,
		&orderproto.UpdateOrderStatusRequest{
			Id: uint32(orderID),
			Status: orderproto.OrderStatus_ORDER_STATUS_PAID,
		},
	)

	return err
}

func (c *Client) CancelOrder(
	ctx context.Context,
	orderID uint,
) error {

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

	_, err := c.client.UpdateOrderStatus(
		ctx,
		&orderproto.UpdateOrderStatusRequest{
			Id: uint32(orderID),
			Status: orderproto.OrderStatus_ORDER_STATUS_CANCELLED,
		},
	)

	return err
}