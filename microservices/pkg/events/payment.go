package events

type PaymentCompleted struct {
	PaymentID uint    `json:"payment_id"`
	OrderID   uint    `json:"order_id"`
	UserID    uint    `json:"user_id"`
	Amount    float64 `json:"amount"`
}