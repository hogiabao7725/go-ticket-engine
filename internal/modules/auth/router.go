package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hogiabao7725/gin-auth-playground/internal/core/config"
	"github.com/hogiabao7725/gin-auth-playground/internal/core/middleware"
	"github.com/hogiabao7725/gin-auth-playground/internal/infra/sqlc"
	infraToken "github.com/hogiabao7725/gin-auth-playground/internal/infra/token"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/features/login"
	get_me "github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/features/me"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/features/register"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/infra"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.RouterGroup, dbPool *pgxpool.Pool, jwtCfg config.JWTConfig) {
	queries := sqlc.New(dbPool)

	// Infrastructure adapters
	userRepo := infra.NewUserRepository(queries)
	passHasher := infra.NewBcryptHasher()
	idGen := infra.NewUUIDGenerator()
	jwtEngine := infraToken.NewJWT(jwtCfg.AccessSecret, jwtCfg.RefreshSecret, jwtCfg.AccessTTL, jwtCfg.RefreshTTL)
	tokenGenerator := infra.NewJWTTokenGenerator(jwtEngine)
	authMiddleware := middleware.NewAuthMiddleware(jwtEngine)

	// Use case handlers
	registerHandler := register.NewHandler(userRepo, passHasher, idGen)
	loginHandler := login.NewHandler(userRepo, passHasher, tokenGenerator)
	getMeHandler := get_me.NewHandler(userRepo)

	// HTTP handlers
	registerHTTPHandler := register.NewHTTPHandler(registerHandler)
	loginHTTPHandler := login.NewHTTPHandler(loginHandler)
	getMeHTTPHandler := get_me.NewHTTPHandler(getMeHandler)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", registerHTTPHandler.Register)
		authGroup.POST("/login", loginHTTPHandler.Login)

		protected := authGroup.Group("/")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/me", getMeHTTPHandler.GetMe)
		}
	}
}
