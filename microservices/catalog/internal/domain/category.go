package domain

import "time"

type Category struct {
	ID           uint
	Name         string
	ParentID     uint
	ImageURL     string
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}