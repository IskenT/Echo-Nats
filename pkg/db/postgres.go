package postgres

import (
	"database/sql"
	"rest_clickhouse/configs"
	"rest_clickhouse/pkg/logger"
)

type DB struct {
	*sql.DB
	log logger.Logger
}

func NewDBConnection(cnf *configs.Config, log logger.Logger) (*DB, error) {
	db, err := sql.Open("postgres", cnf.GetPgDsn())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cnf.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cnf.Postgres.MaxIdleConns)

	return &DB{
		db, log,
	}, nil
}
