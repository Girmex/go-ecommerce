package ports

import "context"

type User struct {
	ID    uint
	Name  string
	Email string
}

type UserClient interface {
	GetUser(
		ctx context.Context,
		userID uint,
	) (*User, error)
}
