package db

import (
	"context"
	"fmt"
	"time"

	"github.com/axosec/vault/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a config, error: %w", err)
	}

	dbConfig.MaxConns = 30
	dbConfig.MinConns = 2
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("error while creating connection to the database!")
	}
	return connPool, nil
}
