package grpc

import (
	"github.com/Girmex/go-ecommerce/microservices/catalog/api/proto"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application"
)

type Handler struct {
	proto.UnimplementedCatalogServiceServer

	service *application.CatalogService
}

func NewHandler(service *application.CatalogService) *Handler {
	return &Handler{
		service: service,
	}
}