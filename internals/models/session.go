package models

import "time"

type Session struct {
	ID        string    `gorm:"primaryKey;size:64"`
	UserID    uint      `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
