package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"
	ports "github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/port"
)

// Mock UserRepository
type mockUserRepository struct {
	users map[string]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {
	user.ID = uint(len(m.users) + 1)
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(
	ctx context.Context,
	id uint,
) (*domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// Mock TokenService
type mockTokenService struct {
	token string
	err   error
}

func (m *mockTokenService) Generate(userID uint) (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.token, nil
}

func (m *mockTokenService) Validate(token string) (uint, error) {
	return 1, nil
}

func TestCreateUser(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{}

	service := NewUserService(
		repository,
		tokenService,
	)

	user, err := service.CreateUser(
		context.Background(),
		"Test User",
		"test@example.com",
		"password123",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Name != "Test User" {
		t.Errorf("expected name Test User, got %s", user.Name)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	if user.PasswordHash == "password123" {
		t.Error("password should be hashed")
	}

	if user.PasswordHash == "" {
		t.Error("password hash should not be empty")
	}
}

func TestCreateUser_InvalidInput(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.CreateUser(
		context.Background(),
		"",
		"test@example.com",
		"password123",
	)

	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf(
			"expected ErrInvalidInput, got %v",
			err,
		)
	}
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.CreateUser(
		context.Background(),
		"Test User",
		"test@example.com",
		"password123",
	)

	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = service.CreateUser(
		context.Background(),
		"Another User",
		"test@example.com",
		"password123",
	)

	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf(
			"expected ErrUserAlreadyExists, got %v",
			err,
		)
	}
}

func TestLogin(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{
		token: "test-token",
	}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.CreateUser(
		context.Background(),
		"Test User",
		"test@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	result, err := service.Login(
		context.Background(),
		"test@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected login result, got nil")
	}

	if result.User == nil {
		t.Fatal("expected user, got nil")
	}

	if result.User.Email != "test@example.com" {
		t.Errorf(
			"expected email test@example.com, got %s",
			result.User.Email,
		)
	}

	if result.Token != "test-token" {
		t.Errorf(
			"expected token test-token, got %s",
			result.Token,
		)
	}
}

func TestLogin_InvalidInput(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{
		token: "test-token",
	}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.Login(
		context.Background(),
		"",
		"password123",
	)

	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf(
			"expected ErrInvalidInput, got %v",
			err,
		)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{
		token: "test-token",
	}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.Login(
		context.Background(),
		"unknown@example.com",
		"password123",
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{
		token: "test-token",
	}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.CreateUser(
		context.Background(),
		"Test User",
		"test@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	_, err = service.Login(
		context.Background(),
		"test@example.com",
		"wrong-password",
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}

func TestLogin_TokenGenerationFails(t *testing.T) {
	tokenErr := errors.New("token generation failed")

	repository := newMockUserRepository()
	tokenService := &mockTokenService{
		err: tokenErr,
	}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.CreateUser(
		context.Background(),
		"Test User",
		"test@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	_, err = service.Login(
		context.Background(),
		"test@example.com",
		"password123",
	)

	if !errors.Is(err, tokenErr) {
		t.Fatalf(
			"expected token generation error, got %v",
			err,
		)
	}
}

func TestGetUser(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{}

	service := NewUserService(
		repository,
		tokenService,
	)

	createdUser, err := service.CreateUser(
		context.Background(),
		"Test User",
		"test@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	user, err := service.GetUser(
		context.Background(),
		createdUser.ID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.ID != createdUser.ID {
		t.Errorf(
			"expected ID %d, got %d",
			createdUser.ID,
			user.ID,
		)
	}

	if user.Email != "test@example.com" {
		t.Errorf(
			"expected email test@example.com, got %s",
			user.Email,
		)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	repository := newMockUserRepository()
	tokenService := &mockTokenService{}

	service := NewUserService(
		repository,
		tokenService,
	)

	_, err := service.GetUser(
		context.Background(),
		999,
	)

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}

// Make sure our mocks implement the interfaces.
var _ ports.UserRepository = (*mockUserRepository)(nil)
var _ ports.TokenService = (*mockTokenService)(nil)
