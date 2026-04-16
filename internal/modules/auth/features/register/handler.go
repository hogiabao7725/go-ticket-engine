package register

import (
	"context"

	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
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
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	name, err := domain.NewName(cmd.Name)
	if err != nil {
		return nil, err
	}

	plainPwd, err := domain.NewPlainPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	hashStr, err := h.hasher.Hash(plainPwd.Value())
	if err != nil {
		return nil, err
	}

	hashedPwd := domain.NewHashedPassword(hashStr)

	userId := h.idGen.Generate()

	user, err := domain.NewUser(userId, name, email, hashedPwd, domain.RoleUser)
	if err != nil {
		return nil, err
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
