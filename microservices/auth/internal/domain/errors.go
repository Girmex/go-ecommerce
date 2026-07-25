package domain

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserAlreadyExists     = errors.New("user already exists")

	//jwt errors
	ErrInvalidToken          = errors.New("invalid token")
	ErrExpiredToken          = errors.New("expired token")
	ErrInvalidTokenType      = errors.New("invalid token type")
	ErrInvalidSigningMethod  = errors.New("invalid signing method")

)
