package bootstrap

import (
	"context"
	"os"

	"github.com/afrizalsebastian/football-team-management/config"
	"github.com/afrizalsebastian/football-team-management/module/database"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FootballManagementApp struct {
	Env    *config.Config
	DBPool *pgxpool.Pool
}

func NewApplication() *FootballManagementApp {
	l := logger.LoggerNew()
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		l.WithContext(ctx).
			Fatal("Failed to get working directory").
			Attr("error", err).
			Msg()
	}

	appConfig, err := config.Init(wd)
	if err != nil {
		l.WithContext(ctx).
			Fatal("Failed to initiate appConfig").
			Attr("error", err).
			Msg()
	}

	pool, err := database.NewDatabasePool(ctx, &database.DatabaseConfig{
		DbUser:     appConfig.DBUser,
		DbPassword: appConfig.DBPassword,
		DbHost:     appConfig.DBHost,
		DbPort:     appConfig.DBPort,
		DbName:     appConfig.DBName,
	})
	if err != nil {
		l.WithContext(ctx).
			Fatal("Failed to initiate db connections").
			Attr("error", err).
			Msg()
	}

	return &FootballManagementApp{
		Env:    appConfig,
		DBPool: pool,
	}
}
