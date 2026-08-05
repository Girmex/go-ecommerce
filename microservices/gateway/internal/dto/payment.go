package dto

type CreatePaymentRequest struct {
	OrderID uint32  `json:"order_id"`
	UserID  uint32  `json:"user_id"`
	Amount  float64 `json:"amount"`
}