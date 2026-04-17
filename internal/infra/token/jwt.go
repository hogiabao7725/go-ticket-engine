package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrEmptyUserID      = errors.New("userID cannot be empty")
	ErrEmptySecret      = errors.New("secret cannot be empty")
	ErrInvalidTTL       = errors.New("TTL must be positive")
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrInvalidToken     = errors.New("invalid token")
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type RefreshClaims struct {
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
	return token.SignedString([]byte(j.accessSecret))
}

func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	if err := validateTokenInput(userID, j.refreshSecret, j.refreshTTL); err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(j.refreshTTL)
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

func (j *JWT) ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.accessSecret), nil
	})
	if err != nil {
		return nil, mapJWTError(err)
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (j *JWT) ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.refreshSecret), nil
	})
	if err != nil {
		return nil, mapJWTError(err)
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

// ValidateAccessToken returns true if the access token is valid.
func (j *JWT) ValidateAccessToken(tokenString string) bool {
	_, err := j.ParseAccessToken(tokenString)
	return err == nil
}

// ValidateRefreshToken returns true if the refresh token is valid.
func (j *JWT) ValidateRefreshToken(tokenString string) bool {
	_, err := j.ParseRefreshToken(tokenString)
	return err == nil
}

// mapJWTError converts standard jwt errors to package errors.
func mapJWTError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return ErrTokenExpired
	case errors.Is(err, jwt.ErrTokenMalformed):
		return ErrTokenMalformed
	case errors.Is(err, jwt.ErrSignatureInvalid):
		return ErrInvalidSignature
	default:
		return err
	}
}

func validateTokenInput(userID, secret string, ttl time.Duration) error {
	if strings.TrimSpace(userID) == "" {
		return ErrEmptyUserID
	}
	if strings.TrimSpace(secret) == "" {
		return ErrEmptySecret
	}
	if ttl <= 0 {
		return ErrInvalidTTL
	}
	return nil
}

// AccessTTL returns the configured access token TTL.
func (j *JWT) AccessTTL() time.Duration {
	return j.accessTTL
}

// RefreshTTL returns the configured refresh token TTL.
func (j *JWT) RefreshTTL() time.Duration {
	return j.refreshTTL
}
