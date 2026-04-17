package domain

import "errors"

var (
	// ID
	ErrEmptyID = errors.New("id cannot be empty")

	// Name
	ErrEmptyName = errors.New("name cannot be empty")

	// Email
	ErrEmptyEmail   = errors.New("email cannot be empty")
	ErrInvalidEmail = errors.New("invalid email format")

	// Password
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrWeakPassword  = errors.New("password is too weak")

	// Role
	ErrInvalidRole = errors.New("invalid role")

	// User
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrUserNotFound      = errors.New("user not found")

	// Authentication
	ErrInvalidCredentials = errors.New("invalid email or password")
)
