package models

import "time"

type ProductModel struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"index"`
	Description string
	CategoryID  uint
	ImageURL    string
	Price       float64
	UserID      uint
	Stock       uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}