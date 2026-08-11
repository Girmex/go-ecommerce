package grpc

import (
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/auth/proto"
)

func toProtoUser(user *domain.User) *proto.User {
	return &proto.User{
		Id:    uint32(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		PhoneVerified: user.PhoneVerified,
		
	}
}
