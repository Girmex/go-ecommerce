package domain

type OrderItem struct {
	ID uint

	OrderID uint

	ProductID uint

	Quantity uint

	UnitPrice float64
}