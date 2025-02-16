package auth

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET_KEY []byte

func init() {
	secretKeyString := os.Getenv("JWT_SECRET_KEY")
	if secretKeyString == "" {
		slog.Error("JWT_SECRET_KEY environment variable not set")
		os.Exit(1)
	}

	JWT_SECRET_KEY = []byte(secretKeyString)
}

func CreateAccessToken(email string) (string, error) {
	return createToken(email, time.Now().Add(15*time.Minute))
}

func CreateRefreshToken(email string) (string, error) {
	return createToken(email, time.Now().Add(7*24*time.Hour))
}

type CustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func createToken(email string, expiration time.Time) (string, error) {
	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	})

	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
