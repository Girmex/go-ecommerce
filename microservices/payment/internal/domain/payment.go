package domain
import "time"


type Payment struct {
	ID uint
	OrderID uint
	UserID uint
	Amount float64
	Status PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}