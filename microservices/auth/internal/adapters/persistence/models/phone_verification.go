package models

import "time"

type PhoneVerificationModel struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"`
	Phone     string    `gorm:"not null"`
	CodeHash  string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
	Used      bool      `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
