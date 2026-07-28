package models
import (
	"time"
)

type OrderModel struct {
	ID uint `gorm:"primaryKey"`

	UserID uint

	Status string

	TotalPrice float64

	Items []OrderItemModel `gorm:"foreignKey:OrderID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}