package repository

import (
	"time"

	"github.com/souvikjs01/auth-cli/internals/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) FindByID(id string) (*models.Session, error) {
	var session models.Session

	err := r.db.
		Preload("User").
		First(&session, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) Delete(id string) error {
	return r.db.
		Delete(&models.Session{}, "id = ?", id).
		Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.db.
		Where("expires_at < ?", time.Now()).
		Delete(&models.Session{}).
		Error
}
