package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github/anurag/altar-be/system/config"
)

func NewPostgresClient(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.PG_URL)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil

}
