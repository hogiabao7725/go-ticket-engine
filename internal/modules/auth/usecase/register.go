package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
)

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterResponse struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

type RegisterUseCase struct {
	repo domain.UserRepository
}

func NewRegisterUseCase(repo domain.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{
		repo: repo,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	user, err := domain.NewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("usecase.Register: %w", err)
	}

	return &RegisterResponse{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Role:      user.Role(),
		CreatedAt: user.CreatedAt(),
	}, nil
}
