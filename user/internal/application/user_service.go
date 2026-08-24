package application

import (
	"context"
	"errors"
	"strings"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/port"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	repository   ports.UserRepository
	tokenService ports.TokenService
}

func NewUserService(
	repository ports.UserRepository,
	tokenService ports.TokenService,
) *UserService {
	return &UserService{
		repository:   repository,
		tokenService: tokenService,
	}
}
type LoginResult struct {
	User  *domain.User
	Token string
}

func (s *UserService) CreateUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*domain.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" || email == "" || password == "" {
		return nil, ErrInvalidInput
	}

	existingUser, err := s.repository.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}





func (s *UserService) Login(
	ctx context.Context,
	email string,
	password string,
) (*LoginResult, error) {
	email = strings.TrimSpace(email)

	if email == "" || password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenService.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}



func (s *UserService) GetUser(
	ctx context.Context,
	id uint,
) (*domain.User, error) {
	return s.repository.GetByID(ctx, id)
}
