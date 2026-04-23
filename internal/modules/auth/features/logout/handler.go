package logout

import (
	"context"
	"fmt"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
)

type Handler struct {
	refreshRepo domain.RefreshTokenRepository
	tokenHasher domain.TokenHasher
}

func NewHandler(refreshRepo domain.RefreshTokenRepository, tokenHasher domain.TokenHasher) *Handler {
	return &Handler{
		refreshRepo: refreshRepo,
		tokenHasher: tokenHasher,
	}
}

func (h *Handler) Execute(ctx context.Context, cmd Command) error {
	if cmd.RefreshToken == "" {
		return nil // nothing to logout
	}

	tokenHash, err := h.tokenHasher.Hash(cmd.RefreshToken)
	if err != nil {
		return fmt.Errorf("auth.features.logout.handler.Execute.Hash: %w", err)
	}

	err = h.refreshRepo.DeleteByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("auth.features.logout.handler.Execute.DeleteByTokenHash: %w", err)
	}

	return nil
}
