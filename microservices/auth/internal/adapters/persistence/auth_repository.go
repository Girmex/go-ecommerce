package persistence

import (
	"context"

	"errors"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/adapters/persistence/models"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

// DeleteProduct implements [ports.CatalogRepository].

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *domain.User) error {

	model := toUserModel(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Copy generated values back to the domain entity
	*user = *toUserDomain(model)

	return nil
}

func (repo *AuthRepository) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {

	var model models.UserModel

	err := repo.db.WithContext(ctx).
		First(&model, id).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return toUserDomain(&model), nil
}

func (repo *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	var model models.UserModel

	err := repo.db.WithContext(ctx).
		Where("email = ?", email).
		First(&model).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return toUserDomain(&model), nil
}
func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, userID uint, refreshToken string,
) error {

	return r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("id = ?", userID).
		Update("refresh_token", refreshToken).
		Error
}
func (r *AuthRepository) MarkPhoneVerified(
	ctx context.Context,
	userID uint,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("id = ?", userID).
		Update("phone_verified", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
