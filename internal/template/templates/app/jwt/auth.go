package jwt

import (
	"log/slog"
	"time"

	"example.com/app"
	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type authService struct {
	secret   string
	duration time.Duration
}

func NewAuthService(secret string, duration time.Duration) app.AuthService {
	return &authService{
		secret:   secret,
		duration: duration,
	}
}

func (s *authService) GenerateToken(userID int) (string, error) {
	claims := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		slog.Error("failed to generate token", "error", err, "userID", userID)
		return "", app.Errorf(app.INTERNAL_ERROR, "failed to generate token")
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, app.Errorf(app.UNAUTHORIZED_ERROR, "invalid token")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return 0, app.Errorf(app.UNAUTHORIZED_ERROR, "invalid token")
	}

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, app.Errorf(app.UNAUTHORIZED_ERROR, "invalid token claims")
}
