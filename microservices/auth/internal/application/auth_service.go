package application

import (
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/ports"
	"golang.org/x/crypto/bcrypt"

	"context"
)

type AuthService struct {
	repository ports.AuthRepository
}

func NewAuthService(repository ports.AuthRepository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	input dto.RegisterInput,
) (*domain.User, error) {

	existingUser, err := s.repository.GetUserByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
	}
	if err := s.repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
