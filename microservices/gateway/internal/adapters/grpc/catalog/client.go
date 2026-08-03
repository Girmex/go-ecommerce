package catalog

import (
	catalogproto "github.com/Girmex/go-ecommerce/microservices/catalog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	catalogproto.CatalogServiceClient
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
		CatalogServiceClient: catalogproto.NewCatalogServiceClient(conn),
	}, nil
}