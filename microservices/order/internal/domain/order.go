package domain
import (
	"time"
)

type Order struct {
    ID         uint
    UserID     uint
    Status     OrderStatus
    TotalPrice float64
    Items []OrderItem
    CreatedAt time.Time
    UpdatedAt time.Time
}




