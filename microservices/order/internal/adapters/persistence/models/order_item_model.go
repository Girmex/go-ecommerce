package models
type OrderItemModel struct {
	ID uint `gorm:"primaryKey"`

	OrderID uint

	ProductID uint

	Quantity uint

	UnitPrice float64
}