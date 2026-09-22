package database

import (
	"context"
	"fmt"

	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

func NewDatabasePool(ctx context.Context, config *DatabaseConfig) (*pgxpool.Pool, error) {
	l := logger.LoggerNew()

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		l.WithContext(ctx).
			Error("Error when initiate pool config").
			Attr("erorr", err).Msg()
		return nil, err
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		l.WithContext(ctx).Error("Error when initiate db pool").Attr("erorr", err).Msg()
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		l.WithContext(ctx).Error("DB Server not respond").Attr("erorr", err).Msg()
		return nil, err
	}

	return pool, nil
}
