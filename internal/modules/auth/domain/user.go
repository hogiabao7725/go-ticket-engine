package domain

import (
	"time"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain/uservo"
)

type User struct {
	id        string
	name      uservo.Name
	email     uservo.Email
	password  uservo.HashedPassword
	role      uservo.Role
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(id string, name uservo.Name, email uservo.Email, hashedPassword uservo.HashedPassword, role uservo.Role) (*User, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	now := time.Now()

	return &User{
		id:        id,
		name:      name,
		email:     email,
		password:  hashedPassword,
		role:      role,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructUser(id string, name uservo.Name, email uservo.Email, hashedPassword uservo.HashedPassword, role uservo.Role, createdAt, updatedAt time.Time) *User {
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

func (u *User) ID() string { return u.id }

func (u *User) Name() uservo.Name { return u.name }

func (u *User) Email() uservo.Email { return u.email }

func (u *User) Role() uservo.Role { return u.role }

func (u *User) PasswordHash() string { return u.password.Value() }

func (u *User) CreatedAt() time.Time { return u.createdAt }

func (u *User) UpdatedAt() time.Time { return u.updatedAt }
