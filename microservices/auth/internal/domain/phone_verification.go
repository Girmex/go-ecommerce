package domain

import "time"

type PhoneVerification struct {
	ID        uint
	UserID    uint
	Phone     string
	CodeHash  string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (v *PhoneVerification) IsExpired(now time.Time) bool {
	return now.After(v.ExpiresAt)
}

func (v *PhoneVerification) MarkUsed() {
	v.Used = true
}
