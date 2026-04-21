package register

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain/uservo"
)

type Handler struct {
	userRepo domain.UserRepository
	hasher   domain.PasswordHasher
	idGen    domain.IdentifierGenerator
}

func NewHandler(userRepo domain.UserRepository, hasher domain.PasswordHasher, idGen domain.IdentifierGenerator) *Handler {
	return &Handler{
		userRepo: userRepo,
		hasher:   hasher,
		idGen:    idGen,
	}
}

func (h *Handler) Execute(ctx context.Context, cmd Command) (*domain.User, error) {
	email, err := uservo.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	name, err := uservo.NewName(cmd.Name)
	if err != nil {
		return nil, err
	}

	plainPwd, err := uservo.NewPlainPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	hashStr, err := h.hasher.Hash(plainPwd.Value())
	if err != nil {
		return nil, err
	}

	hashedPwd := uservo.NewHashedPassword(hashStr)

	userId := h.idGen.Generate()

	user, err := domain.NewUser(userId, name, email, hashedPwd, uservo.RoleUser)
	if err != nil {
		return nil, err
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, err
		}
		return nil, fmt.Errorf("auth.features.register.handler.Execute: %w", err)
	}

	return user, nil
}
