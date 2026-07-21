package grpc

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/catalog/api/proto"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application/dto"
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

func (h *Handler) CreateCategory(
	ctx context.Context,
	req *proto.CreateCategoryRequest,
) (*proto.Category, error) {

	input := dto.CreateCategoryInput{
		Name:         req.Name,
		ParentID:     uint(req.ParentId),
		ImageURL:     req.ImageUrl,
		DisplayOrder: int(req.DisplayOrder),
	}

	category, err := h.service.CreateCategory(ctx, input)
	if err != nil {
		return nil, err
	}
	return toProtoCategory(category), nil
}
