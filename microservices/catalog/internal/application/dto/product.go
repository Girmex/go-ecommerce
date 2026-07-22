package dto

type CreateProductInput struct {
	Name        string
	Description string
	CategoryID  uint32
	ImageURL    string
	Price       float64
	Stock       uint32
}

type UpdateProductInput struct {
	Name        *string
	Description *string
	CategoryID  *uint32
	ImageURL    *string
	Price       *float64
	Stock       *uint32
}

type UpdateStockInput struct {
	Stock uint32
}