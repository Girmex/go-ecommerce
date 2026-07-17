package domain

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrProductNotFound  = errors.New("product not found")
	ErrUnauthorized     = errors.New("unauthorized")

	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrProductAlreadyExists  = errors.New("product already exists")
	ErrProductOutOfStock     = errors.New("product out of stock")
)