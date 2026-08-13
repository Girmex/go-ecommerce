package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")

	//jwt errors
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("expired token")
	ErrInvalidTokenType     = errors.New("invalid token type")
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	//phone verification errors
	ErrPhoneVerificationNotFound = errors.New("phone verification not found")
	ErrPhoneVerificationExpired  = errors.New("phone verification expired")
	ErrPhoneVerificationUsed     = errors.New("phone verification already used")
	ErrInvalidPhoneVerification  = errors.New("invalid phone verification code")
)
