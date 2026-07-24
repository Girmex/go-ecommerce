package grpc

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/auth/proto"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/dto"
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
