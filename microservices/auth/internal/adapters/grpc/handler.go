package grpc

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/auth/proto"
)

type Handler struct {
	proto.UnimplementedAuthServiceServer

	service *application.AuthService
}

func NewHandler(service *application.AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(
	ctx context.Context,
	req *proto.RegisterRequest,
) (*proto.RegisterResponse, error) {

	input := dto.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.service.Register(ctx, input)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &proto.RegisterResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *Handler) Login(
	ctx context.Context,
	req *proto.LoginRequest,
) (*proto.LoginResponse, error) {

	input := dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.service.Login(ctx, input)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &proto.LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		User:         toProtoUser(output.User),
	}, nil
}

func (h *Handler) GetUser(
	ctx context.Context,
	req *proto.GetUserRequest,
) (*proto.User, error) {

	user, err := h.service.GetUserByID(ctx, uint(req.Id))
	if err != nil {
		return nil, toStatusError(err)
	}
	return toProtoUser(user), nil
}
