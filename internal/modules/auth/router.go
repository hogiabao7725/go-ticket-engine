package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hogiabao7725/go-ticket-engine/internal/core/config"
	"github.com/hogiabao7725/go-ticket-engine/internal/infra/sqlc"
	infraToken "github.com/hogiabao7725/go-ticket-engine/internal/infra/token"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/features/login"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/features/register"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/infra"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.RouterGroup, dbPool *pgxpool.Pool, jwtCfg config.JWTConfig) {
	queries := sqlc.New(dbPool)

	// Infrastructure adapters
	userRepo := infra.NewUserRepository(queries)
	passHasher := infra.NewBcryptHasher()
	idGen := infra.NewUUIDGenerator()
	tokenGenerator := infra.NewJWTTokenGenerator(
		infraToken.NewJWT(jwtCfg.AccessSecret, jwtCfg.RefreshSecret, jwtCfg.AccessTTL, jwtCfg.RefreshTTL),
	)

	// Use case handlers
	registerHandler := register.NewHandler(userRepo, passHasher, idGen)
	loginHandler := login.NewHandler(userRepo, passHasher, tokenGenerator)

	// HTTP handlers
	registerHTTPHandler := register.NewHTTPHandler(registerHandler)
	loginHTTPHandler := login.NewHTTPHandler(loginHandler)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", registerHTTPHandler.Register)
		authGroup.POST("/login", loginHTTPHandler.Login)
	}
}
