package grpc

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/catalog/api/proto"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/application/dto"
	"google.golang.org/protobuf/types/known/emptypb"
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
		ParentID:     req.ParentId,
		ImageURL:     req.ImageUrl,
		DisplayOrder: req.DisplayOrder,
	}

	category, err := h.service.CreateCategory(ctx, input)
	if err != nil {
		return nil, toStatusError(err)

	}
	return toProtoCategory(category), nil
}

func (h *Handler) GetCategory(
	ctx context.Context,
	req *proto.GetCategoryRequest,
) (*proto.Category, error) {

	category, err := h.service.GetCategoryByID(ctx, uint(req.Id))

	if err != nil {
		return nil, toStatusError(err)

	}
	return toProtoCategory(category), nil
}

func (h *Handler) ListCategories(
	ctx context.Context, req *emptypb.Empty,
) (*proto.ListCategoriesResponse, error) {

	categories, err := h.service.GetCategories(ctx)
	if err != nil {
		return nil, toStatusError(err)
	}
	response := make([]*proto.Category, 0, len(categories))
	for _, category := range categories {
		response = append(response, toProtoCategory(category))
	}
	return &proto.ListCategoriesResponse{
		Categories: response,
	}, nil
}
func (h *Handler) UpdateCategory(
	ctx context.Context,
	req *proto.UpdateCategoryRequest,
) (*proto.Category, error) {
	input := dto.UpdateCategoryInput{
	Name:         req.Name,
	ParentID:     req.ParentId,
	ImageURL:     req.ImageUrl,
	DisplayOrder: req.DisplayOrder,
}
	category, err := h.service.UpdateCategory(ctx, uint(req.Id), input)
	if err != nil {
		return nil, toStatusError(err)
	}
	return toProtoCategory(category), nil
}

func (h *Handler) DeleteCategory(
	ctx context.Context,
	req *proto.GetCategoryRequest,
) (*emptypb.Empty, error) {

	if err := h.service.DeleteCategory(ctx, uint(req.Id)); err != nil {
		return nil, toStatusError(err)
	}
	return &emptypb.Empty{}, nil
}
