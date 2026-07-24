package grpc

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/catalog/proto"
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
func (h *Handler) CreateProduct(
	ctx context.Context,
	req *proto.CreateProductRequest,
) (*proto.Product, error) {

	const userID uint32 = 1
	input := dto.CreateProductInput{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryId,
		ImageURL:    req.ImageUrl,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	product, err := h.service.CreateProduct(
		ctx,
		input,
		userID,
	)
	if err != nil {
		return nil, toStatusError(err)
	}
	return toProtoProduct(product), nil
}

func (h *Handler) GetProduct(
	ctx context.Context,
	req *proto.GetProductRequest,
) (*proto.Product, error) {

	product, err := h.service.GetProductByID(ctx, uint(req.Id))
	if err != nil {
		return nil, toStatusError(err)
	}
	return toProtoProduct(product), nil
}

func (h *Handler) GetProducts(
	ctx context.Context, req *emptypb.Empty,
) (*proto.ListProductsResponse, error) {

	products, err := h.service.GetProducts(ctx)
	if err != nil {
		return nil, toStatusError(err)
	}
	response := make([]*proto.Product, 0, len(products))
	for _, product := range products {
		response = append(response, toProtoProduct(product))
	}
	return &proto.ListProductsResponse{
		Products: response,
	}, nil
}

func (h *Handler) UpdateProduct(
	ctx context.Context,
	req *proto.UpdateProductRequest,
) (*proto.Product, error) {

	const userID uint32 = 1

	input := dto.UpdateProductInput{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryId,
		ImageURL:    req.ImageUrl,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	product, err := h.service.UpdateProduct(ctx, uint(req.Id), input, userID)
	if err != nil {
		return nil, toStatusError(err)
	}
	return toProtoProduct(product), nil
}

func (h *Handler) DeleteProduct(
	ctx context.Context,
	req *proto.GetProductRequest,
) (*emptypb.Empty, error) {

	const userID uint32 = 1
	if err := h.service.DeleteProduct(ctx, uint(req.Id), userID); err != nil {
		return nil, toStatusError(err)
	}
	return &emptypb.Empty{}, nil
}
func (h *Handler) GetSellerProducts(
	ctx context.Context,
	req *proto.GetSellerProductsRequest,
) (*proto.ListProductsResponse, error) {

	const userID uint = 1
	products, err := h.service.GetSellerProducts(ctx, userID)
	if err != nil {
		return nil, toStatusError(err)
	}
	response := make([]*proto.Product, 0, len(products))
	for _, product := range products {
		response = append(response, toProtoProduct(product))
	}
	return &proto.ListProductsResponse{
		Products: response,
	}, nil
}
func (h *Handler) UpdateProductStock(
	ctx context.Context,
	req *proto.UpdateStockRequest,
) (*proto.Product, error) {
	const userID uint32 = 1
	input := dto.UpdateStockInput{
		Stock: req.Stock,
	}
	product, err := h.service.UpdateProductStock(ctx, uint(req.Id), input, userID)
	if err != nil {
		return nil, toStatusError(err)
	}
	return toProtoProduct(product), nil
}
