package grpc

import (
	"github.com/Girmex/go-ecommerce/microservices/auth/api/proto"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
)

func toProtoUser(user *domain.User) *proto.User {
	return &proto.User{
		Id:    uint32(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}
}

func toDomainUser(user *proto.User) *domain.User {
	return &domain.User{
		ID:    uint(user.Id),
		Name:  user.Name,
		Email: user.Email,
	}
}

