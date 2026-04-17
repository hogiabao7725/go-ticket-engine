package domain

import "time"

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) error
}

type IdentifierGenerator interface {
	Generate() string
}

type TokenResult struct {
	Token     string
	ExpiresIn time.Duration
}

type TokenGenerator interface {
	GenerateAccessToken(userID, role string) (TokenResult, error)
	GenerateRefreshToken(userID string) (string, error)
}
