package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidTokenInput = errors.New("invalid token input parameters")
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}

func validateTokenInput(userID, secret string, ttl time.Duration) error {
	if strings.TrimSpace(userID) == "" {
		return ErrInvalidTokenInput
	}
	if strings.TrimSpace(secret) == "" {
		return ErrInvalidTokenInput
	}
	if ttl <= 0 {
		return ErrInvalidTokenInput
	}
	return nil
}

func GenerateAccessToken(userID string, role string, secretKey string, ttl time.Duration) (string, error) {
	if err := validateTokenInput(userID, secretKey, ttl); err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(ttl)
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken(userID, secret string, ttl time.Duration) (string, error) {
	if err := validateTokenInput(userID, secret, ttl); err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(ttl)
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
