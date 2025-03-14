package jwt

import (
	"time"

	"example.com/rest/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	secret   string
	duration time.Duration
}

func NewAuthService(secret string, duration time.Duration) *AuthService {
	return &AuthService{
		secret:   secret,
		duration: duration,
	}
}

type customClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(userID int) (string, error) {
	claims := customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", domain.Errorf(domain.INTERNAL_ERROR, "failed to generate token").Wrap(err)
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	// parse the token
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		// This function mainly needs to return the key for validating the token.
		// We verify the signing method to prevent attacks with other signing methods.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid token")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return 0, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid token")
	}

	// get our custom claims
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid token")
}
