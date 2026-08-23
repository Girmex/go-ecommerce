package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	TokenKey  contextKey = "token"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				writeUnauthorized(w, "authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") ||
				parts[1] == "" {
				writeUnauthorized(w, "invalid authorization header")
				return
			}

			tokenString := parts[1]

			token, err := jwt.Parse(
				tokenString,
				func(token *jwt.Token) (interface{}, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, jwt.ErrSignatureInvalid
					}

					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeUnauthorized(w, "invalid token claims")
				return
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				writeUnauthorized(w, "invalid user id")
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				uint(userID),
			)

			ctx = context.WithValue(
				ctx,
				TokenKey,
				tokenString,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

func UserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

func Token(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(TokenKey).(string)
	return token, ok
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_, _ = w.Write([]byte(
		`{"code":"UNAUTHORIZED","message":"` + message + `"}`,
	))
}
