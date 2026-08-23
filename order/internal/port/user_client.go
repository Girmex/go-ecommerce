package port

import "context"

type UserInfo struct {
	ID    string
	Name  string
	Email string
}

type UserClient interface {
	GetUser(ctx context.Context, userID string) (*UserInfo, error)
}