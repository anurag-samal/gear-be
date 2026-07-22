package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github/anurag/altar-be/system/config"
)

type RedisConfig struct {
	Client       *redis.Client
	isConnected  bool
}

func NewRedisClient(ctx context.Context, cfg *config.Config) (*RedisConfig, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Redis.Address,
		Password:        cfg.Redis.Password,
		DB:              0,
		MinIdleConns:    2,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
		PoolTimeout:     4 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	slog.Info("redis client configured",
		"addr", cfg.Redis.Address,
		"min_idle", 2,
		"max_idle", 10,
	)

	return &RedisConfig{Client: client, isConnected: true}, nil
}

func (r *RedisConfig) HealthCheck(ctx context.Context) bool {
	if r.Client == nil {
		return false
	}
	if err := r.Client.Ping(ctx).Err(); err != nil {
		slog.Error("redis health check failed", "error", err)
		r.isConnected = false
		return false
	}
	return true
}

func (r *RedisConfig) IsConnected() bool {
	return r.isConnected
}

func (r *RedisConfig) Disconnect() {
	if r.Client == nil {
		return
	}
	slog.Info("disconnecting redis")
	if err := r.Client.Close(); err != nil {
		slog.Error("redis close error", "error", err)
	}
	r.isConnected = false
	slog.Info("redis disconnected")
}
