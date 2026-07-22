package domain

import "time"

type Product struct {
	ID          uint
	Name        string
	Description string
	CategoryID  uint32
	ImageURL    string
	Price       float64
	UserID      uint32
	Stock       uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}