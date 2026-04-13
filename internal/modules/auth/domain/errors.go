package domain

import (
	"errors"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidName        = errors.New("invalid name (must not be empty)")
	ErrWeakPassword       = errors.New("password too weak (min 6 characters)")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
