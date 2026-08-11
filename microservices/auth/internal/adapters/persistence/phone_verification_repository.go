package persistence

import (
	"context"
	"errors"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/ports"
	"gorm.io/gorm"
)

var _ ports.PhoneVerificationRepository = (*PhoneVerificationRepository)(nil)

type PhoneVerificationRepository struct {
	db *gorm.DB
}

func NewPhoneVerificationRepository(db *gorm.DB) *PhoneVerificationRepository {
	return &PhoneVerificationRepository{
		db: db,
	}
}

func (r *PhoneVerificationRepository) Create(
	ctx context.Context,
	verification *domain.PhoneVerification,
) error {
	model := toPhoneVerificationModel(verification)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Copy generated values, such as ID, back to domain.
	*verification = *toPhoneVerificationDomain(model)

	return nil
}

func (r *PhoneVerificationRepository) GetLatestByUserID(
	ctx context.Context,
	userID uint,
) (*domain.PhoneVerification, error) {
	var model models.PhoneVerificationModel

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPhoneVerificationNotFound
	}

	if err != nil {
		return nil, err
	}

	return toPhoneVerificationDomain(&model), nil
}

func (r *PhoneVerificationRepository) MarkUsed(
	ctx context.Context,
	verification *domain.PhoneVerification,
) error {
	verification.MarkUsed()

	result := r.db.WithContext(ctx).
		Model(&models.PhoneVerificationModel{}).
		Where("id = ?", verification.ID).
		Updates(map[string]interface{}{
			"used":       verification.Used,
			"updated_at": verification.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrPhoneVerificationNotFound
	}

	return nil
}
