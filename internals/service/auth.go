package service

import (
	"errors"
	"strings"
	"time"

	"github.com/souvikjs01/auth-cli/internals/models"
	repository "github.com/souvikjs01/auth-cli/internals/repositories"
	"github.com/souvikjs01/auth-cli/internals/security"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountLocked      = errors.New("account is temporarily locked")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if name == "" || email == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *AuthService) Login(input LoginInput) (*models.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if email == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, ErrAccountLocked
	}

	if !security.VerifyPassword(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()

	user.FailedAttempts = 0
	user.LockedUntil = nil
	user.LastLoginAt = &now

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
