package refresh

import (
	"context"
	"fmt"
	"time"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
)

type Result struct {
	AccessToken      string
	ExpiresIn        int64
	RefreshToken     string
	RefreshExpiresIn int64
}

type Handler struct {
	jwtTokenGenerator domain.TokenGenerator
	refreshRepo       domain.RefreshTokenRepository
	tokenHasher       domain.TokenHasher
	userRepo          domain.UserRepository
	idGen             domain.IdentifierGenerator
}

func NewHandler(
	jwtTokenGenerator domain.TokenGenerator,
	refreshRepo domain.RefreshTokenRepository,
	tokenHasher domain.TokenHasher,
	userRepo domain.UserRepository,
	idGen domain.IdentifierGenerator,
) *Handler {
	return &Handler{
		jwtTokenGenerator: jwtTokenGenerator,
		refreshRepo:       refreshRepo,
		tokenHasher:       tokenHasher,
		userRepo:          userRepo,
		idGen:             idGen,
	}
}

func (h *Handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	if cmd.RefreshToken == "" {
		return nil, domain.ErrMissingToken
	}

	// 1. Hash the incoming refresh token
	tokenHash, err := h.tokenHasher.Hash(cmd.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.Hash: %w", err)
	}

	// 2. Find the token in the DB
	token, err := h.refreshRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.FindByTokenHash: %w", err)
	}
	if token == nil {
		return nil, domain.ErrInvalidToken
	}

	// 3. Check if expired
	if time.Now().After(token.ExpiresAt()) {
		return nil, domain.ErrTokenExpired
	}

	// 4. Revoke the token by deleting it from the database (Rotation)
	err = h.refreshRepo.DeleteByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.DeleteByTokenHash: %w", err)
	}

	// 5. Get the user to generate a new Access Token (we need the Role)
	user, err := h.userRepo.FindByID(ctx, token.UserID())
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.FindUserByID: %w", err)
	}

	// 6. Generate new Access Token and Refresh Token
	tokenResult, err := h.jwtTokenGenerator.GenerateAccessToken(user.ID(), user.Role().String())
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.GenerateAccessToken: %w", err)
	}

	newRefreshTokenResult, err := h.jwtTokenGenerator.GenerateRefreshToken(user.ID())
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.GenerateRefreshToken: %w", err)
	}

	// 7. Hash the new refresh token and store it
	newTokenHash, err := h.tokenHasher.Hash(newRefreshTokenResult.Token)
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.HashNewToken: %w", err)
	}

	newID := h.idGen.Generate()
	now := time.Now()
	expiresAt := now.Add(newRefreshTokenResult.ExpiresIn) 
	
	newSession := domain.NewRefreshToken(newID, user.ID(), newTokenHash, expiresAt, now)
	err = h.refreshRepo.Create(ctx, newSession)
	if err != nil {
		return nil, fmt.Errorf("auth.features.refresh.handler.Execute.CreateSession: %w", err)
	}

	return &Result{
		AccessToken:      tokenResult.Token,
		ExpiresIn:        int64(tokenResult.ExpiresIn.Seconds()),
		RefreshToken:     newRefreshTokenResult.Token,
		RefreshExpiresIn: int64(newRefreshTokenResult.ExpiresIn.Seconds()),
	}, nil
}
