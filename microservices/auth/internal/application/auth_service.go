package application

import (
	"errors"
	"strconv"
	"time"

	"context"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/helpers"
	"github.com/Girmex/go-ecommerce/microservices/auth/internal/ports"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
	"github.com/Girmex/go-ecommerce/microservices/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository            ports.AuthRepository
	eventPublisher        ports.EventPublisher
	phoneVerificationRepo ports.PhoneVerificationRepository
	jwtManager            *jwt.JWTManager
}

func NewAuthService(repository ports.AuthRepository, eventPublisher ports.EventPublisher, phoneVerificationRepo ports.PhoneVerificationRepository, jwtManager *jwt.JWTManager,
) *AuthService {
	return &AuthService{
		repository:            repository,
		eventPublisher:        eventPublisher,
		phoneVerificationRepo: phoneVerificationRepo,
		jwtManager:            jwtManager,
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
		Name:          input.Name,
		Email:         input.Email,
		PasswordHash:  string(hash),
		Phone:         input.Phone,
		PhoneVerified: false,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	input dto.LoginInput,
) (*dto.LoginOutput, error) {

	user, err := s.repository.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {

		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	if err := s.repository.UpdateRefreshToken(
		ctx,
		user.ID,
		refreshToken,
	); err != nil {
		return nil, err
	}
	return &dto.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
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

func (s *AuthService) RequestPhoneVerification(
	ctx context.Context,
	user *domain.User,
) error {
	code, err := helpers.GenerateOTP()
	if err != nil {
		return err
	}

	codeHash, err := bcrypt.GenerateFromPassword(
		[]byte(code),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	verification := &domain.PhoneVerification{
		UserID:    user.ID,
		Phone:     user.Phone,
		CodeHash:  string(codeHash),
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}

	if err := s.phoneVerificationRepo.Create(
		ctx,
		verification,
	); err != nil {
		return err
	}

	event := events.UserPhoneVerification{
		UserID: user.ID,
		Phone:  user.Phone,
		Code:   code,
	}

	if err := s.eventPublisher.Publish(
		ctx,
		"user.phone_verification",
		strconv.FormatUint(uint64(user.ID), 10),
		event,
	); err != nil {
		return err
	}

	return nil
}
