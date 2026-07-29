package domain

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrForbidden     = errors.New("forbidden")
)