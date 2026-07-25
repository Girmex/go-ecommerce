package middleware

import "context"

type contextKey string

const userIDKey contextKey = "user_id"

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserID(ctx context.Context) (uint32, bool) {
	userID, ok := ctx.Value(userIDKey).(uint32)
	return userID, ok
}
