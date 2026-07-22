package grpc

import (
	"github.com/Girmex/go-ecommerce/microservices/catalog/api/proto"
	"github.com/Girmex/go-ecommerce/microservices/catalog/internal/domain"
)

func toProtoCategory(category *domain.Category) *proto.Category {
	return &proto.Category{
		Id:           uint32(category.ID),
		Name:         category.Name,
		ParentId:     uint32(category.ParentID),
		ImageUrl:     category.ImageURL,
		DisplayOrder: int32(category.DisplayOrder),
	}
}

func toProtoProduct(product *domain.Product) *proto.Product {
	return &proto.Product{
		Id:          uint32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		CategoryId:  uint32(product.CategoryID),
		ImageUrl:    product.ImageURL,
		Price:       product.Price,
		UserId:      uint32(product.UserID),
		Stock:       uint32(product.Stock),
	}
}