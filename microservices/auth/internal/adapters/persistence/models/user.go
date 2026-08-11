package models

import "time"

type UserModel struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Email         string `gorm:"uniqueIndex"`
	PasswordHash  string
	RefreshToken  string `gorm:"type:text"`
	Phone         string `gorm:"uniqueIndex"`
	PhoneVerified bool   `gorm:"not null;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
