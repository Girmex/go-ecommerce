package middleware

import (
	"context"
	"strings"

	"github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]bool{
	"/catalog.v1.CatalogService/ListProducts": true,
	"/catalog.v1.CatalogService/GetProduct":   true,
}

func AuthInterceptor(
	jwtManager *jwt.JWTManager,
) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Skip authentication for public endpoints
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(
				codes.Unauthenticated,
				"missing metadata",
			)
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(
				codes.Unauthenticated,
				"missing authorization header",
			)
		}

		authHeader := authHeaders[0]

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Error(
				codes.Unauthenticated,
				"invalid authorization header",
			)
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			return nil, status.Error(
				codes.Unauthenticated,
				"invalid token",
			)
		}

		ctx = WithUserID(ctx, uint(claims.UserID))

		return handler(ctx, req)
	}
}
