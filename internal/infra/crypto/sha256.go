package crypto

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
)

type SHA256 struct{}

func (s SHA256) Hash(plain string) (string, error) {
	trimmed := strings.TrimSpace(plain)
	if trimmed == "" {
		return "", ErrEmptyInput
	}
	sum := sha256.Sum256([]byte(trimmed))
	return hex.EncodeToString(sum[:]), nil
}

func (s SHA256) Compare(hashed, token string) error {
	trimmedHashed := strings.TrimSpace(hashed)
	if trimmedHashed == "" {
		return ErrInvalidHash
	}

	calculated, err := s.Hash(token)
	if err != nil {
		return fmt.Errorf("crypto: sha256 hash calculation failed: %w", err)
	}

	if subtle.ConstantTimeCompare([]byte(trimmedHashed), []byte(calculated)) != 1 {
		return ErrMismatched
	}
	return nil
}
