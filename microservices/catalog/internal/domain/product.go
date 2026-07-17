package domain

import "time"

type Product struct {
	ID          uint
	Name        string
	Description string
	CategoryID  uint
	ImageURL    string
	Price       float64
	UserID      uint
	Stock       uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}