package domain

type User struct {
	ID            uint
	Name          string
	Email         string
	Phone         string
	PasswordHash  string
	PhoneVerified bool
	RefreshToken  string
}
