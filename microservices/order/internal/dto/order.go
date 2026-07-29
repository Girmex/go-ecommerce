package dto

type CreateOrderItemInput struct {
	ProductID uint32
	Quantity  uint32
}

type CreateOrderInput struct {
	Items []CreateOrderItemInput
}