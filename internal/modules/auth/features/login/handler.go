package login

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain/uservo"
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
	refreshRepo    domain.RefreshTokenRepository
	tokenHasher    domain.TokenHasher
	idGen          domain.IdentifierGenerator
}

func NewHandler(
	userRepo domain.UserRepository,
	passwordHasher domain.PasswordHasher,
	tokenGen domain.TokenGenerator,
	refreshRepo domain.RefreshTokenRepository,
	tokenHasher domain.TokenHasher,
	idGen domain.IdentifierGenerator,
) *Handler {
	return &Handler{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGen:       tokenGen,
		refreshRepo:    refreshRepo,
		tokenHasher:    tokenHasher,
		idGen:          idGen,
	}
}

// Execute returns Result and the RefreshToken's TokenResult so HTTP layer can set the cookie
func (h *Handler) Execute(ctx context.Context, cmd Command) (*Result, domain.TokenResult, error) {
	email, err := uservo.NewEmail(cmd.Email)
	if err != nil {
		return nil, domain.TokenResult{}, domain.ErrInvalidCredentials
	}

	pass, err := uservo.NewPlainPassword(cmd.Password)
	if err != nil {
		return nil, domain.TokenResult{}, domain.ErrInvalidCredentials
	}

	user, err := h.userRepo.FindByEmail(ctx, email.String())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.TokenResult{}, domain.ErrInvalidCredentials
		}
		return nil, domain.TokenResult{}, fmt.Errorf("auth.features.login.handler.Execute: %w", err)
	}

	// compare password
	if err := h.passwordHasher.Compare(user.PasswordHash(), pass.Value()); err != nil {
		return nil, domain.TokenResult{}, domain.ErrInvalidCredentials
	}

	userId := user.ID()
	role := user.Role().String()

	// generate tokens
	tokenResult, err := h.tokenGen.GenerateAccessToken(userId, role)
	if err != nil {
		return nil, domain.TokenResult{}, fmt.Errorf("auth.features.login.handler.Execute.GenerateAccessToken: %w", err)
	}

	// generate refresh token
	refreshTokenResult, err := h.tokenGen.GenerateRefreshToken(userId)
	if err != nil {
		return nil, domain.TokenResult{}, fmt.Errorf("auth.features.login.handler.Execute.GenerateRefreshToken: %w", err)
	}

	// hash the refresh token before saving
	tokenHash, err := h.tokenHasher.Hash(refreshTokenResult.Token)
	if err != nil {
		return nil, domain.TokenResult{}, fmt.Errorf("auth.features.login.handler.Execute.HashToken: %w", err)
	}

	// save to DB
	id := h.idGen.Generate()
	now := time.Now()
	expiresAt := now.Add(refreshTokenResult.ExpiresIn)
	
	newSession := domain.NewRefreshToken(id, userId, tokenHash, expiresAt, now)
	err = h.refreshRepo.Create(ctx, newSession)
	if err != nil {
		return nil, domain.TokenResult{}, fmt.Errorf("auth.features.login.handler.Execute.CreateRefreshToken: %w", err)
	}

	return &Result{
		AccessToken: tokenResult.Token,
		ExpiresIn:   int64(tokenResult.ExpiresIn.Seconds()),
		User:        user,
	}, refreshTokenResult, nil
}
