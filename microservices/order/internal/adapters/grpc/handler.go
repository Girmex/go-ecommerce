package grpc

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/grpc/middleware"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/dto"
	proto "github.com/Girmex/go-ecommerce/microservices/order/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	proto.UnimplementedOrderServiceServer

	service *application.OrderService
}

func NewHandler(service *application.OrderService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateOrder(
	ctx context.Context,
	req *proto.CreateOrderRequest,
) (*proto.Order, error) {

	userID, ok := middleware.UserID(ctx)
	if !ok {
		return nil, status.Error(
			codes.Unauthenticated,
			"user id missing",
		)
	}

	items := make([]dto.CreateOrderItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, dto.CreateOrderItemInput{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	input := dto.CreateOrderInput{
		Items: items,
	}

	order, err := h.service.CreateOrder(
		ctx,
		input,
		uint32(userID),
	)
	if err != nil {
		return nil, toStatusError(err)
	}

	return toProtoOrder(order), nil
}

func (h *Handler) GetOrder(
	ctx context.Context,
	req *proto.GetOrderRequest,
) (*proto.Order, error) {

	order, err := h.service.GetOrder(
		ctx,
		uint(req.Id),
	)
	if err != nil {
		return nil, toStatusError(err)
	}

	return toProtoOrder(order), nil
}

func (h *Handler) ListOrders(
	ctx context.Context,
	_ *emptypb.Empty,
) (*proto.ListOrdersResponse, error) {

	userID, ok := middleware.UserID(ctx)
	if !ok {
		return nil, status.Error(
			codes.Unauthenticated,
			"user id missing",
		)
	}

	orders, err := h.service.GetOrdersByUser(
		ctx,
		uint32(userID),
	)
	if err != nil {
		return nil, toStatusError(err)
	}

	response := make([]*proto.Order, 0, len(orders))

	for _, order := range orders {
		response = append(
			response,
			toProtoOrder(&order),
		)
	}

	return &proto.ListOrdersResponse{
		Orders: response,
	}, nil
}

func (h *Handler) CancelOrder(
	ctx context.Context,
	req *proto.GetOrderRequest,
) (*emptypb.Empty, error) {

	userID, ok := middleware.UserID(ctx)
	if !ok {
		return nil, status.Error(
			codes.Unauthenticated,
			"user id missing",
		)
	}

	err := h.service.CancelOrder(
		ctx,
		uint(req.Id),
		uint32(userID),
	)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &emptypb.Empty{}, nil
}