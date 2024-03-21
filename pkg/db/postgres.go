package postgres

import (
	"context"
	"fmt"
	"rest_clickhouse/configs"
	"rest_clickhouse/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
	log logger.Logger
}

func NewDBConnection(ctx context.Context, cnf *configs.Config, log logger.Logger) (*DB, error) {
	db, err := pgxpool.New(ctx, cnf.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("error create new db: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error ping db: %w", err)
	}

	return &DB{
		db, log,
	}, nil
}
