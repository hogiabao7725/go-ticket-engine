package usecase

import (
	"context"
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

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return &RegisterResponse{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Role:      user.Role(),
		CreatedAt: user.CreatedAt(),
	}, nil
}
