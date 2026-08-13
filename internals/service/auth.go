package service

import (
	"errors"
	"strings"

	"github.com/souvikjs01/auth-cli/internals/models"
	repository "github.com/souvikjs01/auth-cli/internals/repositories"
	"github.com/souvikjs01/auth-cli/internals/security"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidInput       = errors.New("invalid input")
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
