package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey"`
	Name           string    `gorm:"size:100;not null"`
	Email          string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash   string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	LastLoginAt    *time.Time
	FailedAttempts int `gorm:"not null;default:0"`
	LockedUntil    *time.Time
	MFAEnabled     bool `gorm:"not null;default:false"`
	MFASecret      *string
}
