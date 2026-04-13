package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hogiabao7725/go-ticket-engine/internal/core/config"
	"github.com/hogiabao7725/go-ticket-engine/internal/core/database"
	"github.com/hogiabao7725/go-ticket-engine/internal/core/middleware"
	"github.com/hogiabao7725/go-ticket-engine/internal/infra/sqlc"
	authDelivery "github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/delivery/http"
	authRepository "github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/repository"
	authUsecase "github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/usecase"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/health"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	cfg := loadConfig()
	initLogger(cfg.Server.Env)

	pgPool := connectPostgres(ctx, cfg)
	defer pgPool.Close()

	rd := connectRedis(ctx, cfg)
	defer rd.Close()

	router := setupRouter()
	registerRoutes(router, pgPool)
	startServer(router, cfg)
}

// ======= Config & Logger ======= //

func loadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	return cfg
}

func initLogger(env string) {
	zerolog.TimeFieldFormat = time.RFC3339
	if strings.EqualFold(env, "development") {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
		return
	}

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// ======= Database & Redis ======= //
func connectPostgres(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	pool, err := database.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize postgres")
	}
	log.Info().Msg("successfully connected to postgres")
	return pool
}

func connectRedis(ctx context.Context, cfg *config.Config) *redis.Client {
	client, err := database.NewRedisClient(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize redis")
	}
	log.Info().Msg("successfully connected to redis")
	return client
}

// ======= Router & Server ======= //
func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	return r
}

func registerRoutes(r *gin.Engine, dbPool *pgxpool.Pool) {
	v1 := r.Group("/api/v1")

	// Health check
	health.NewHealthHandler().RegisterRoutes(v1)

	// Auth module
	queries := sqlc.New(dbPool)
	userRepo := authRepository.NewUserRepository(queries)
	registerUC := authUsecase.NewRegisterUseCase(userRepo)
	authHandler := authDelivery.NewAuthHandler(registerUC)
	authHandler.RegisterRoutes(v1)

	// Log all registered routes
	for _, route := range r.Routes() {
		log.Info().
			Str("method", route.Method).
			Str("path", route.Path).
			Msg("route registered")
	}
}

func startServer(r *gin.Engine, cfg *config.Config) {
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Info().Str("addr", addr).Msg("starting server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server error")
	}
}
