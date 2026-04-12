package domain

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hogiabao7725/go-ticket-engine/pkg/hash"
)

type User struct {
	id        string
	name      string
	email     string
	password  string // hashed
	role      string // user, organizer, admin
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(name, email, plainPassword string) (*User, error) {
	normalizedName, err := normalizeAndValidateName(name)
	if err != nil {
		return nil, err
	}

	normalizedEmail, err := normalizeAndValidateEmail(email)
	if err != nil {
		return nil, err
	}

	if err := validatePassword(plainPassword); err != nil {
		return nil, err
	}

	hashedPassword, err := hash.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		id:        uuid.New().String(),
		name:      normalizedName,
		email:     normalizedEmail,
		password:  hashedPassword,
		role:      "user",
		createdAt: now,
		updatedAt: now,
	}, nil
}

func normalizeAndValidateName(name string) (string, error) {
	normalizedName := strings.Join(strings.Fields(name), " ")
	if normalizedName == "" {
		return "", ErrInvalidName
	}
	return normalizedName, nil
}

func normalizeAndValidateEmail(email string) (string, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	if _, err := mail.ParseAddress(normalizedEmail); err != nil {
		return "", ErrInvalidEmail
	}
	return normalizedEmail, nil
}

func validatePassword(plainPassword string) error {
	if len(plainPassword) < 6 {
		return ErrWeakPassword
	}
	return nil
}

func ReconstructUser(id, name, email, hashedPassword, role string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		password:  hashedPassword,
		role:      role,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) VerifyPassword(plain string) bool {
	return hash.ComparePassword(u.password, plain) == nil
}

func (u *User) ID() string { return u.id }

func (u *User) Name() string { return u.name }

func (u *User) Email() string { return u.email }

func (u *User) Role() string { return u.role }

func (u *User) PasswordHash() string { return u.password }

func (u *User) CreatedAt() time.Time { return u.createdAt }

func (u *User) UpdatedAt() time.Time { return u.updatedAt }
