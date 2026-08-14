// Package service implements the core business logic for authentication,
// session management, and TOTP-based multi-factor authentication.
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
	ErrUsernameAlreadyExists = errors.New("username already registered")
	ErrInvalidInput          = errors.New("invalid input")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrAccountLocked         = errors.New("account is temporarily locked")
)

// AuthService handles user registration and login, including account
// lockout after repeated failed attempts.
type AuthService struct {
	userRepo         *repository.UserRepository
	maxLoginAttempts int
	lockoutDuration  time.Duration
}

// NewAuthService creates an AuthService with configurable lockout settings.
func NewAuthService(userRepo *repository.UserRepository, maxLoginAttempts int, lockoutDuration time.Duration) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		maxLoginAttempts: maxLoginAttempts,
		lockoutDuration:  lockoutDuration,
	}
}

// RegisterInput holds the fields required to register a new user.
type RegisterInput struct {
	Username string
	Password string
}

// LoginInput holds the fields required to authenticate a user.
type LoginInput struct {
	Username string
	Password string
}

// Register creates a new user account. It validates input, checks for
// duplicate usernames, and stores the password as a bcrypt hash.
func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	username := strings.TrimSpace(input.Username)

	if username == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	exists, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUsernameAlreadyExists
	}

	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user by username and password. It enforces account
// lockout after maxLoginAttempts consecutive failures and records the last
// successful login time.
func (s *AuthService) Login(input LoginInput) (*models.User, error) {
	username := strings.TrimSpace(input.Username)

	if username == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	// Check whether the account is currently locked.
	if user.LockedUntil != nil {
		if time.Now().Before(*user.LockedUntil) {
			return nil, ErrAccountLocked
		}

		// Lock has expired — reset the counter.
		user.LockedUntil = nil
		user.FailedAttempts = 0

		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	}

	// Verify password.
	if !security.VerifyPassword(
		input.Password,
		user.PasswordHash,
	) {
		user.FailedAttempts++

		// Lock the account if too many failed attempts.
		if user.FailedAttempts >= s.maxLoginAttempts {
			lockedUntil := time.Now().Add(s.lockoutDuration)
			user.LockedUntil = &lockedUntil
		}

		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}

		return nil, ErrInvalidCredentials
	}

	// Successful login — reset failed attempts and record timestamp.
	now := time.Now()

	user.FailedAttempts = 0
	user.LockedUntil = nil
	user.LastLoginAt = &now

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
