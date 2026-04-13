package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/go-ticket-engine/internal/infra/sqlc"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type userRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queriers *sqlc.Queries) domain.UserRepository {
	return &userRepository{
		queries: queriers,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.queries.CreateUser(ctx, fromDomainUser(user))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("repository.Create: %w", err)
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("repository.FindByEmail: %w", err)
	}
	return toDomainUser(dbUser), nil
}

// toDomainUser: sqlc.User -> domain.User
func toDomainUser(dbUser sqlc.User) *domain.User {
	return domain.ReconstructUser(
		dbUser.ID,
		dbUser.Name,
		dbUser.Email,
		dbUser.Password,
		dbUser.Role,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)
}

// fromDomainUser: domain.User -> sqlc.CreateUserParams
func fromDomainUser(user *domain.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Password:  user.PasswordHash(),
		Role:      user.Role(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}
