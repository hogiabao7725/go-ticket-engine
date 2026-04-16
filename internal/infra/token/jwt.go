package token

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidTokenInput = errors.New("invalid token input parameters")

type accessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type refreshClaims struct {
	jwt.RegisteredClaims
}

type JWT struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWT(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWT {
	return &JWT{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (j *JWT) GenerateAccessToken(userID, role string) (string, error) {
	if err := validateTokenInput(userID, j.accessSecret, j.accessTTL); err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(j.accessTTL)
	claims := accessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecret))
}

func GenerateRefreshToken(userID, secret string, ttl time.Duration) (string, error) {
	if err := validateTokenInput(userID, secret, ttl); err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(ttl)
	claims := refreshClaims{
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
