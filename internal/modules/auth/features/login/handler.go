package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
)

type Result struct {
	AccessToken string
	ExpiresIn   int64
	User        *domain.User
}

type Handler struct {
	userRepo       domain.UserRepository
	passwordHasher domain.PasswordHasher
	tokenGen       domain.TokenGenerator
}

func NewHandler(userRepo domain.UserRepository, passwordHasher domain.PasswordHasher, tokenGen domain.TokenGenerator) *Handler {
	return &Handler{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGen:       tokenGen,
	}
}

func (h *Handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	pass, err := domain.NewPlainPassword(cmd.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := h.userRepo.FindByEmail(ctx, email.String())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("auth.features.login.handler.Execute: %w", err)
	}

	// compare password
	if err := h.passwordHasher.Compare(user.PasswordHash(), pass.Value()); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	userId := user.ID()
	role := user.Role().String()

	// generate tokens
	tokenResult, err := h.tokenGen.GenerateAccessToken(userId, role)
	if err != nil {
		return nil, fmt.Errorf("auth.features.login.handler.Execute: %w", err)
	}

	/*
		TODO: Generate and store refresh token if needed.
		For now, we only generate access token since refresh token handling is not implemented yet.
	*/

	return &Result{
		AccessToken: tokenResult.Token,
		ExpiresIn:   int64(tokenResult.ExpiresIn.Seconds()),
		User:        user,
	}, nil

}
