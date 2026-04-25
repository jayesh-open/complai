package testutil

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresContainer(ctx context.Context) (connString string, cleanup func(), err error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("complai_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		return "", nil, fmt.Errorf("testutil: start postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return "", nil, fmt.Errorf("testutil: get postgres connection string: %w", err)
	}

	cleanup = func() {
		_ = pgContainer.Terminate(ctx)
	}

	return connStr, cleanup, nil
}

func NewRedisContainer(ctx context.Context) (connString string, cleanup func(), err error) {
	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections"),
		),
	)
	if err != nil {
		return "", nil, fmt.Errorf("testutil: start redis container: %w", err)
	}

	host, err := redisContainer.Host(ctx)
	if err != nil {
		_ = redisContainer.Terminate(ctx)
		return "", nil, fmt.Errorf("testutil: get redis host: %w", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		_ = redisContainer.Terminate(ctx)
		return "", nil, fmt.Errorf("testutil: get redis port: %w", err)
	}

	connStr := fmt.Sprintf("%s:%s", host, port.Port())

	cleanup = func() {
		_ = redisContainer.Terminate(ctx)
	}

	return connStr, cleanup, nil
}
