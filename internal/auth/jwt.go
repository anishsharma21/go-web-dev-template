package auth

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
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

func CreateAccessToken(id int) (string, error) {
	return createToken(id, time.Now().Add(15*time.Minute))
}

func CreateRefreshToken(id int) (string, error) {
	return createToken(id, time.Now().Add(7*24*time.Hour))
}

type CustomClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func createToken(id int, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		id,
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

func VerifyToken(tokenString string) (int, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v\n", t.Header["alg"])
		}

		return JWT_SECRET_KEY, nil
	})
	if err != nil {
		return -1, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return -1, fmt.Errorf("invalid token claims")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return -1, fmt.Errorf("invalid token claims: id not found")
	}

	return int(idFloat), nil
}
