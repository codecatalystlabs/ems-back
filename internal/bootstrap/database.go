package bootstrap

import (
	"context"
	"dispatch/internal/platform/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	pcfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, err
	}
	pcfg.MaxConns = int32(cfg.MaxOpenConns)
	pcfg.MinConns = 2
	pcfg.MaxConnIdleTime = 5 * time.Minute
	pcfg.MaxConnLifetime = cfg.ConnMaxLifetime
	return pgxpool.NewWithConfig(ctx, pcfg)
}
