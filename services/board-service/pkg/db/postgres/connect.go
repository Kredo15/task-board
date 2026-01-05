package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kredo15/task-board/services/board-service/internal/config"
)

type Client struct {
	Pool *pgxpool.Pool
}

type DatabaseInterface interface {
	GetPool() *pgxpool.Pool
	Close()
	Ping(ctx context.Context) error
}

func NewClient(cfg *config.Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dsn := cfg.Postgres.GetDSN()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = int32(cfg.Postgres.MaxConns)
	config.MinConns = int32(cfg.Postgres.MinConns)
	config.MaxConnLifetime = cfg.Postgres.MaxConnLifetime
	config.MaxConnIdleTime = cfg.Postgres.MaxConnIdleTime
	config.HealthCheckPeriod = cfg.Postgres.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{Pool: pool}, nil
}

func (c *Client) GetPool() *pgxpool.Pool {
	return c.Pool
}

func (c *Client) Close() {
	c.Pool.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}
