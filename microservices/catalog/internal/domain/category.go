package domain

import "time"

type Category struct {
	ID           uint
	Name         string
	ParentID     uint32
	ImageURL     string
	DisplayOrder int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}