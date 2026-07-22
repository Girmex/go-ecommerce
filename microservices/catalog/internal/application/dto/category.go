package dto

type CreateCategoryInput struct {
	Name         string
	ParentID     uint32
	ImageURL     string
	DisplayOrder int32
}

type UpdateCategoryInput struct {
	Name         *string
	ParentID     *uint32
	ImageURL     *string
	DisplayOrder *int32
}