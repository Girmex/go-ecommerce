package models

import "time"

type CategoryModel struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"index"`
	ParentID     uint32
	ImageURL     string
	DisplayOrder int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}