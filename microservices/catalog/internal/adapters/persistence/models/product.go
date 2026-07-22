package models

import "time"

type ProductModel struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"index"`
	Description string
	CategoryID  uint32
	ImageURL    string
	Price       float64
	UserID      uint32
	Stock       uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}