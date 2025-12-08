package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Config for connection pooling
type Config struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration // optional
	// Retry settings
	RetryMaxWait time.Duration // total time to keep retrying
}

// Open opens DB, sets pool config, and pings with retry.
func Open(dsn string, cfg Config) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("empty dsn")
	}

	// sensible defaults
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 5 * time.Minute
	}
	if cfg.RetryMaxWait == 0 {
		cfg.RetryMaxWait = 30 * time.Second
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Apply pooling config
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	if cfg.ConnMaxIdleTime > 0 {
		// available since Go 1.15+
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	// Ping with exponential backoff + jitter up to RetryMaxWait
	deadline := time.Now().Add(cfg.RetryMaxWait)
	attempt := 0
	for {
		attempt++
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			return db, nil
		}

		if time.Now().After(deadline) {
			_ = db.Close()
			return nil, fmt.Errorf("database unreachable after %s: %w", cfg.RetryMaxWait, err)
		}

		// backoff: base * 2^attempt + jitter
		base := 250 * time.Millisecond
		sleep := base * (1 << uint(min(attempt, 6))) // cap exponent
		// jitter +-25%
		jitter := time.Duration(float64(sleep) * (rand.Float64()*0.5 - 0.25))
		sleep = sleep + jitter
		if sleep < 100*time.Millisecond {
			sleep = 100 * time.Millisecond
		}
		time.Sleep(sleep)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
