package events

type UserPhoneVerification struct {
	UserID uint   `json:"user_id"`
	Phone  string `json:"phone"`
	Code   string `json:"code"`
}
