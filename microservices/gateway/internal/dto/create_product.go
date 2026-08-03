package dto

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  uint32  `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Stock       uint32  `json:"stock"`
}