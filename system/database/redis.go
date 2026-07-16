package database

import (
	"context"
	"github/anurag/altar-be/system/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:     cfg.RDS_ADDRESS,
		Password: cfg.RDS_PASSWORD,
		DB:       0,
	})

	if err := rds.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rds, nil
}
