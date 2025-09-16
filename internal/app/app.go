package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tamboto2000/otaqku-tasks/internal/config"
	"github.com/tamboto2000/otaqku-tasks/internal/database"
	"github.com/vinovest/sqlx"
)

type App struct {
	cfg     config.Config
	db      *sqlx.DB
	httpSrv *http.Server
}

func NewApp(cfg config.Config) (*App, error) {
	// Initiate database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}

	logger := createLogger(cfg.Logging)

	// Repositories
	repos := newRepositories(db, logger)

	// Services
	svcs := newServices(cfg, repos, logger)

	// Register HTTP handlers
	router := echo.New()
	registerHandlers(router, svcs)

	// HTTP server
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPServer.Port),
		Handler: router,
	}

	return &App{
		cfg:     cfg,
		db:      db,
		httpSrv: httpSrv,
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

func (a *App) RunHTTPServer() error {
	return a.httpSrv.ListenAndServe()
}

func (a *App) Shutdown() error {
	if err := a.db.Close(); err != nil {
		return err
	}

	if err := a.httpSrv.Shutdown(context.Background()); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
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
