package ports

import (
	"context"

	authproto "github.com/Girmex/go-ecommerce/microservices/auth/proto"
)

type AuthClient interface {
	GetUser(ctx context.Context, userID uint) (*authproto.User, error)
}
