package models

import "time"

type CategoryModel struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"index"`
	ParentID     uint
	ImageURL     string
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}