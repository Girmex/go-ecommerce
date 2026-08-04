package dto

type CreateCategoryRequest struct {
	Name         string `json:"name"`
	ParentID     uint32 `json:"parent_id"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int32  `json:"display_order"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty"`
	ParentID     *uint32 `json:"parent_id,omitempty"`
	ImageURL     *string `json:"image_url,omitempty"`
	DisplayOrder *int32  `json:"display_order,omitempty"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  uint32  `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Stock       uint32  `json:"stock"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	CategoryID  *uint32  `json:"category_id,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *uint32  `json:"stock,omitempty"`
}

type UpdateStockRequest struct {
	Stock uint32 `json:"stock"`
}
