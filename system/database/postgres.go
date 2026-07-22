package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github/anurag/altar-be/system/config"
)

type PostgresConfig struct {
	Pool         *pgxpool.Pool
	isConnected  bool
}

func NewPostgresClient(ctx context.Context, cfg *config.Config) (*PostgresConfig, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	slog.Info("postgres pool configured",
		"min_conns", poolCfg.MinConns,
		"max_conns", poolCfg.MaxConns,
		"max_conn_lifetime", poolCfg.MaxConnLifetime,
		"max_conn_idle_time", poolCfg.MaxConnIdleTime,
	)

	return &PostgresConfig{Pool: pool, isConnected: true}, nil
}

func (p *PostgresConfig) HealthCheck(ctx context.Context) bool {
	if p.Pool == nil {
		return false
	}
	if err := p.Pool.Ping(ctx); err != nil {
		slog.Error("postgres health check failed", "error", err)
		p.isConnected = false
		return false
	}
	return true
}

func (p *PostgresConfig) IsConnected() bool {
	return p.isConnected
}

func (p *PostgresConfig) Disconnect() {
	if p.Pool == nil {
		return
	}
	slog.Info("disconnecting postgres")
	p.Pool.Close()
	p.isConnected = false
	slog.Info("postgres disconnected")
}
