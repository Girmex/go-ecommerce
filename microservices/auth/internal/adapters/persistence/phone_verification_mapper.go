package persistence

import (
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
)

func toPhoneVerificationModel(
	verification *domain.PhoneVerification,
) *models.PhoneVerificationModel {
	return &models.PhoneVerificationModel{
		ID:        verification.ID,
		UserID:    verification.UserID,
		Phone:     verification.Phone,
		CodeHash:  verification.CodeHash,
		ExpiresAt: verification.ExpiresAt,
		Used:      verification.Used,
		CreatedAt: verification.CreatedAt,
		UpdatedAt: verification.UpdatedAt,
	}
}

func toPhoneVerificationDomain(
	model *models.PhoneVerificationModel,
) *domain.PhoneVerification {
	return &domain.PhoneVerification{
		ID:        model.ID,
		UserID:    model.UserID,
		Phone:     model.Phone,
		CodeHash:  model.CodeHash,
		ExpiresAt: model.ExpiresAt,
		Used:      model.Used,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}