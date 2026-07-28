package domain

type OrderStatus string

const (
	OrderPending         OrderStatus = "PENDING"
	OrderAwaitingPayment OrderStatus = "AWAITING_PAYMENT"
	OrderPaid            OrderStatus = "PAID"
	OrderCancelled       OrderStatus = "CANCELLED"
	OrderShipped         OrderStatus = "SHIPPED"
	OrderDelivered       OrderStatus = "DELIVERED"
)