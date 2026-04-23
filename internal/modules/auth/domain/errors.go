package domain

import (
	"errors"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain/uservo"
)

var (
	// ID
	ErrEmptyID = errors.New("id cannot be empty")

	// User
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrUserNotFound      = errors.New("user not found")

	// Authentication
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Tokens & Session
	ErrInvalidToken = errors.New("invalid refresh token")
	ErrTokenExpired = errors.New("refresh token has expired")
	ErrTokenRevoked = errors.New("refresh token has been revoked")
	ErrMissingToken = errors.New("refresh token is missing")
)

// Validation errors (Aliased from VO)
var (
	ErrEmptyName     = uservo.ErrEmptyName
	ErrEmptyEmail    = uservo.ErrEmptyEmail
	ErrInvalidEmail  = uservo.ErrInvalidEmail
	ErrEmptyPassword = uservo.ErrEmptyPassword
	ErrWeakPassword  = uservo.ErrWeakPassword
	ErrInvalidRole   = uservo.ErrInvalidRole
)
