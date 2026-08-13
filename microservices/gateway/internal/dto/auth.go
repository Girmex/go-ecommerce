package dto

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyPhoneRequest struct {
	UserID uint32 `json:"user_id"`
	Code   string `json:"code"`
}
