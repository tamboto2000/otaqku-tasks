package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/tamboto2000/otaqku-tasks/internal/config"
	"github.com/tamboto2000/otaqku-tasks/internal/database"
	"github.com/vinovest/sqlx"
)

type App struct {
	cfg config.Config
	db  *sqlx.DB
}

func NewApp(cfg config.Config) (*App, error) {
	// Initiate database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg: cfg,
		db:  db,
	}, nil
}

func (a *App) RunDatabaseMigration(ctx context.Context) error {
	// Run database migration
	err := database.RunMigration(context.Background(), a.cfg.Database, a.db)
	if err != nil {
		return err
	}

	return nil
}

func createLogger(logCfg config.Logging) *slog.Logger {
	logHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: logCfg.WithSource,
		Level:     logCfg.Level,
	})

	logger := slog.New(logHandler)

	return logger
}
