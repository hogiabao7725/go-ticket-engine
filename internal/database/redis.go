package database

import (
	"context"
	"fmt"

	"github.com/hogiabao7725/go-ticket-engine/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {

	pingCtx, cancel := context.WithTimeout(ctx, cfg.Redis.ConnectTimeout)
	defer cancel()

	// Config Redis
	rd := redis.NewClient(&redis.Options{
		Addr:        cfg.Redis.Addr(),
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.DB,
		DialTimeout: cfg.Redis.ConnectTimeout,
	})

	if err := rd.Ping(pingCtx).Err(); err != nil {
		_ = rd.Close()
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", cfg.Redis.Addr(), err)
	}

	return rd, nil
}
