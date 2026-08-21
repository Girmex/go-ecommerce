package application

import (
	"context"
	"errors"
	"strings"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/ports"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidInput      = errors.New("invalid input")
)

type UserService struct {
	repository ports.UserRepository
}

func NewUserService(repository ports.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
) (*domain.User, error) {

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" || email == "" || passwordHash == "" {
		return nil, ErrInvalidInput
	}

	existingUser, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(
	ctx context.Context,
	id uint,
) (*domain.User, error) {

	return s.repository.GetByID(ctx, id)
}
