package dto

type CreateCategoryInput struct {
	Name         string
	ParentID     uint
	ImageURL     string
	DisplayOrder int
}

type UpdateCategoryInput struct {
	Name         *string
	ParentID     *uint
	ImageURL     *string
	DisplayOrder *int
}