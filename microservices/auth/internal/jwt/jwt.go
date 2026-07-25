package jwt

import (
	"time"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
}

type Claims struct {
	UserID uint
	Email  string
	TokenType string

	jwtv5.RegisteredClaims
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
	}
}

func (j *JWTManager) GenerateAccessToken(
	user *domain.User,
) (string, error) {

	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		TokenType: "access",
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
			IssuedAt: jwtv5.NewNumericDate(time.Now()),
			NotBefore: jwtv5.NewNumericDate(time.Now()),
		},
	}

	token := jwtv5.NewWithClaims(
		jwtv5.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secret)
}

func (j *JWTManager) GenerateRefreshToken(
	user *domain.User,
) (string, error) {

	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		TokenType: "refresh",
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(
				time.Now().Add(7 * 24 * time.Hour),
			),
			IssuedAt: jwtv5.NewNumericDate(time.Now()),
			NotBefore: jwtv5.NewNumericDate(time.Now()),
		},
	}

	token := jwtv5.NewWithClaims(
		jwtv5.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secret)
}