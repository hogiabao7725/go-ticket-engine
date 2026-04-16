package crypto

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func (Bcrypt) Hash(plain string) (string, error) {
	trimmed := strings.TrimSpace(plain)
	if trimmed == "" {
		return "", ErrEmptyInput
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(trimmed), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("crypto: bcrypt hash failed: %w", err)
	}
	return string(bytes), nil
}

func (Bcrypt) Compare(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrMismatched
		}
		return fmt.Errorf("crypto: bcrypt compare failed: %w", err)
	}
	return nil
}
