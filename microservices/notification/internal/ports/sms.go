package ports

import "context"

type SMSSender interface {
	Send(ctx context.Context, phone string, message string) error
}
