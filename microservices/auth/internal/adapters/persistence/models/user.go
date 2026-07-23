package models
import "time"

type UserModel struct {
    ID           uint   `gorm:"primaryKey"`
    Name         string
    Email        string `gorm:"uniqueIndex"`
    PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}