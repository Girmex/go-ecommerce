package cataloggrpc

import (
	"context"

	catalogproto "github.com/Girmex/go-ecommerce/microservices/catalog/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client catalogproto.CatalogServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		client: catalogproto.NewCatalogServiceClient(conn),
	}
}

func (c *Client) GetProduct(
	ctx context.Context,
	productID uint,
) (*catalogproto.Product, error) {

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

	return c.client.GetProduct(
		ctx,
		&catalogproto.GetProductRequest{
			Id: uint32(productID),
		},
	)
}