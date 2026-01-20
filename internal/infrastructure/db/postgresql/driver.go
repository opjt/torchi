package postgresql

import (
	"context"
	"fmt"
	"time"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Database struct {
	*pgxpool.Pool
}

// NewPoolConfig는 pgxpool.Config를 생성합니다.
func NewPoolConfig(connectionURL string) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Connection Pool 설정
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	return poolConfig, nil
}

// NewDatabase는 즉시 연결을 수행하고 lifecycle에 종료 훅만 등록합니다.
func NewDatabase(lc fx.Lifecycle, env config.Env, logger *log.Logger) (*Database, error) {
	logger.Info("connecting to database...")

	// 1. Pool 설정 생성
	poolConfig, err := NewPoolConfig(env.DB.URL)
	if err != nil {
		return nil, err
	}

	// 2. 즉시 연결 수행
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connected successfully")

	db := &Database{Pool: pool}

	// 3. 종료 시에만 훅 사용
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing database connection...")
			db.Pool.Close()
			logger.Info("database connection closed")
			return nil
		},
	})

	return db, nil
}

// NewQueries는 sqlc Queries를 생성합니다.
func NewQueries(db *Database) *Queries {
	return New(db.Pool)
}

var Module = fx.Module("postgresql",
	fx.Provide(
		NewDatabase,
		NewQueries,
	),
)
