package events

type UserVerification struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}