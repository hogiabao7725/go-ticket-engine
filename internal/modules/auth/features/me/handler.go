package get_me

import (
	"context"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
)

type UserDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type Handler struct {
	userRepo domain.UserRepository
}

func NewHandler(userRepo domain.UserRepository) *Handler {
	return &Handler{userRepo: userRepo}
}

func (h *Handler) Execute(ctx context.Context, cmd Command) (UserDTO, error) {
	user, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return UserDTO{}, err
	}

	return UserDTO{
		ID:        user.ID(),
		Name:      user.Name().String(),
		Email:     user.Email().String(),
		Role:      user.Role().String(),
		CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
