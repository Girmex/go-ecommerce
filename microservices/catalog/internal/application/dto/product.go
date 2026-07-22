package dto

type CreateProductInput struct {
	Name        string
	Description string
	CategoryID  uint
	ImageURL    string
	Price       float64
	Stock       int32
}

type UpdateProductInput struct {
	Name        *string
	Description *string
	CategoryID  *uint32
	ImageURL    *string
	Price       *float64
	Stock       *int32
}

type UpdateStockInput struct {
	Stock int32
}