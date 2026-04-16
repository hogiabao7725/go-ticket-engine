package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hogiabao7725/go-ticket-engine/internal/infra/sqlc"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/features/register"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/infra"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.RouterGroup, dbPool *pgxpool.Pool) {
	queries := sqlc.New(dbPool)

	// Infrastructure adapters
	userRepo := infra.NewUserRepository(queries)
	hasher := infra.NewBcryptHasher()
	idGen := infra.NewUUIDGenerator()

	// Use case handlers
	registerHandler := register.NewHandler(userRepo, hasher, idGen)

	// HTTP handlers
	registerHTTPHandler := register.NewHTTPHandler(registerHandler)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", registerHTTPHandler.Register)
	}
}
