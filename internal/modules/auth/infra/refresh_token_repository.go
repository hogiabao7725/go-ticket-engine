package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/gin-auth-playground/internal/infra/sqlc"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
	"github.com/jackc/pgx/v5"
)

type RefreshTokenRepository struct {
	queries *sqlc.Queries
}

var _ domain.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

func NewRefreshTokenRepository(queries *sqlc.Queries) *RefreshTokenRepository {
	return &RefreshTokenRepository{queries: queries}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	params := sqlc.CreateRefreshTokenParams{
		ID:         token.ID(),
		UserID:     token.UserID(),
		TokenHash:  token.TokenHash(),
		ExpiresAt:  token.ExpiresAt(),
		CreatedAt:  token.CreatedAt(),
	}
	err := r.queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return fmt.Errorf("auth.infra.refresh_token_repo.Create: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	dbToken, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Return nil, nil when not found is standard
		}
		return nil, fmt.Errorf("auth.infra.refresh_token_repo.FindByTokenHash: %w", err)
	}

	return domain.NewRefreshToken(
		dbToken.ID,
		dbToken.UserID,
		dbToken.TokenHash,
		dbToken.ExpiresAt,
		dbToken.CreatedAt,
	), nil
}

func (r *RefreshTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	err := r.queries.DeleteRefreshTokenByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("auth.infra.refresh_token_repo.DeleteByTokenHash: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	err := r.queries.DeleteRefreshTokensByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("auth.infra.refresh_token_repo.DeleteByUserID: %w", err)
	}
	return nil
}
