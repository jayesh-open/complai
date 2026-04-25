package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

var DefaultPoolConfig = PoolConfig{
	MaxConns:          25,
	MinConns:          5,
	MaxConnLifetime:   30 * time.Minute,
	MaxConnIdleTime:   5 * time.Minute,
	HealthCheckPeriod: 30 * time.Second,
}

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	return NewPoolWithConfig(ctx, connString, DefaultPoolConfig)
}

func NewPoolWithConfig(ctx context.Context, connString string, cfg PoolConfig) (*pgxpool.Pool, error) {
	if connString == "" {
		connString = os.Getenv("DATABASE_URL")
	}
	if connString == "" {
		return nil, fmt.Errorf("db: connection string is required (set DATABASE_URL or pass directly)")
	}

	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("db: parse connection string: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("db: create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return pool, nil
}
