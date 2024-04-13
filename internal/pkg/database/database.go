package database

import (
	"banner/internal/pkg/logger"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const driverName = "pgx"

func CreateDBConnection(url string, log *zap.Logger) (*sqlx.DB, error) {

	pgCfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("can't parse pg config: %w", err)
	}

	pgLog := logger.NewPgxLogger(log)
	pgCfg.Tracer = pgLog

	nativeDB := stdlib.OpenDB(*pgCfg)
	if err != nil {
		return nil, err
	}

	nativeDB.SetMaxOpenConns(10)
	nativeDB.SetMaxIdleConns(5)

	return sqlx.NewDb(nativeDB, driverName), nil

}
