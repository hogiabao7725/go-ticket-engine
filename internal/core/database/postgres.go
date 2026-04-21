package database

import (
	"context"
	"fmt"

	"github.com/hogiabao7725/gin-auth-playground/internal/core/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(ctx, cfg.DB.ConnectTimeout)
	defer cancel()

	// Set up connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.DB.DSN())
	if err != nil {
		return nil, fmt.Errorf("error parsing db config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.DB.MaxConns)
	poolConfig.MinConns = int32(cfg.DB.MinConns)
	poolConfig.MaxConnLifetime = cfg.DB.ConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DB.ConnIdleTime

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return pool, nil
}
