package http

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100" example:"Test User"`
	Email    string `json:"email" validate:"required,email" example:"test@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"testpwd123"`
}

type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Test User"`
	Email string `json:"email" example:"test@example.com"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"test@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"testpwd123"`
}

type LoginResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Test User"`
	Email string `json:"email" example:"test@example.com"`
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}
