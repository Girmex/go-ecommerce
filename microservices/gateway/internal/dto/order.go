package dto

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID uint32 `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}
type UpdateOrderStatusRequest struct {
	Status int32 `json:"status"`
}
