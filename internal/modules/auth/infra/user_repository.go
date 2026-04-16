package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/go-ticket-engine/internal/infra/sqlc"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

// Ensure UserRepository implements domain.UserRepository
var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	params := sqlc.CreateUserParams{
		ID:        user.ID(),
		Name:      user.Name().String(),
		Email:     user.Email().String(),
		Password:  user.PasswordHash(),
		Role:      user.Role().String(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
	err := r.queries.CreateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("infra.user_repo.Create: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
