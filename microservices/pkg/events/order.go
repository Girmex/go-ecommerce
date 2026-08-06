package events

type OrderCreated struct {
	OrderID uint    `json:"order_id"`
	UserID  uint    `json:"user_id"`
	Amount  float64 `json:"amount"`
}