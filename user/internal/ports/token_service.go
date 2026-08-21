package ports

type TokenService interface {
	Generate(userID uint) (string, error)
	Validate(token string) (uint, error)
}
