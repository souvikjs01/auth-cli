package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/souvikjs01/auth-cli/internals/models"
	repository "github.com/souvikjs01/auth-cli/internals/repositories"
)

var ErrSessionExpired = errors.New("session has expired")

type SessionService struct {
	sessionRepo    *repository.SessionRepository
	sessionTimeout time.Duration
}

func NewSessionService(
	sessionRepo *repository.SessionRepository,
	sessionTimeout time.Duration,
) *SessionService {
	return &SessionService{
		sessionRepo:    sessionRepo,
		sessionTimeout: sessionTimeout,
	}
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s *SessionService) CreateSession(
	userID uint,
) (*models.Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTimeout),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetValidSession(
	sessionID string,
) (*models.Session, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.Delete(sessionID)

		return nil, ErrSessionExpired
	}

	return session, nil
}

func (s *SessionService) Logout(sessionID string) error {
	return s.sessionRepo.Delete(sessionID)
}
