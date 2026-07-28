package persistence

import (
	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/adapters/persistence/models"
)

func toOrderModel(order *domain.Order) *models.OrderModel {
	return &models.OrderModel{
		ID:         order.ID,
		UserID:     order.UserID,
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
		Items:      toOrderItemModels(order.Items),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func toOrderDomain(model *models.OrderModel) *domain.Order {
	return &domain.Order{
		ID:         model.ID,
		UserID:     model.UserID,
		Status:     domain.OrderStatus(model.Status),
		TotalPrice: model.TotalPrice,
		Items:      toOrderItemDomains(model.Items),
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func toOrderItemModel(item *domain.OrderItem) *models.OrderItemModel {
	return &models.OrderItemModel{
		ID:        item.ID,
		OrderID:   item.OrderID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
	}
}

func toOrderItemDomain(model *models.OrderItemModel) *domain.OrderItem {
	return &domain.OrderItem{
		ID:        model.ID,
		OrderID:   model.OrderID,
		ProductID: model.ProductID,
		Quantity:  model.Quantity,
		UnitPrice: model.UnitPrice,
	}
}

func toOrderItemModels(items []domain.OrderItem) []models.OrderItemModel {
	result := make([]models.OrderItemModel, 0, len(items))

	for _, item := range items {
		result = append(result, *toOrderItemModel(&item))
	}

	return result
}

func toOrderItemDomains(items []models.OrderItemModel) []domain.OrderItem {
	result := make([]domain.OrderItem, 0, len(items))

	for _, item := range items {
		result = append(result, *toOrderItemDomain(&item))
	}

	return result
}