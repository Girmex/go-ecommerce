package dto
import "github.com/Girmex/go-ecommerce/microservices/auth/internal/domain"

type RegisterInput struct {
    Name     string
    Email    string
    Password string
}

type LoginInput struct {
    Email    string
    Password string
}

type LoginOutput struct {
    AccessToken  string
    RefreshToken string
    User         *domain.User
}