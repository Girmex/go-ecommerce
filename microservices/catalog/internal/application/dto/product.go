package dto

type CreateProductInput struct {
	Name        string
	Description string
	CategoryID  uint
	ImageURL    string
	Price       float64
	Stock       int
}

type UpdateProductInput struct {
	Name        *string
	Description *string
	CategoryID  *uint
	ImageURL    *string
	Price       *float64
	Stock       *int
}

type UpdateStockInput struct {
	Stock int
}