package hash

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordEmpty      = errors.New("password cannot be empty")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func HashPassword(password string) (string, error) {

	valPassword, err := validatePassword(password)
	if err != nil {
		return "", err
	}

	hashedByte, err := bcrypt.GenerateFromPassword([]byte(valPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedByte), nil
}

func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("failed to compare password: %w", err)
	}
	return nil
}

func validatePassword(password string) (string, error) {
	valPassword := strings.TrimSpace(password)
	if valPassword == "" {
		return "", ErrPasswordEmpty
	}
	return valPassword, nil
}
