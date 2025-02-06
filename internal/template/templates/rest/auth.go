package rest

type AuthService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (int, error)
}
