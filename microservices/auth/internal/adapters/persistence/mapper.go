package persistence

import (
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
)

func toUserModel(user *domain.User) *models.UserModel {
	return &models.UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		PasswordHash:  user.PasswordHash,
		RefreshToken:  user.RefreshToken,
	}
}

func toUserDomain(model *models.UserModel) *domain.User {
	return &domain.User{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		PasswordHash:  model.PasswordHash,
		RefreshToken:  model.RefreshToken,

	}
}
