package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Role         string
}