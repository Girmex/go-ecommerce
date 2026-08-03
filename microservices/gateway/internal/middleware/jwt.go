package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtpkg "github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func JWTMiddleware(jwtManager *jwtpkg.JWTManager) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Claims(ctx context.Context) (*jwtpkg.Claims, bool) {

	claims, ok := ctx.Value(ClaimsKey).(*jwtpkg.Claims)

	return claims, ok
}
