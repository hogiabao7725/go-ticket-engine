package http

import (
	"context"

	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/usecase"
)

type (
	Registerer interface {
		Execute(ctx context.Context, req usecase.RegisterRequest) (*usecase.RegisterResponse, error)
	}
)
