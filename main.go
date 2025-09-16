package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tamboto2000/otaqku-tasks/internal/app"
	"github.com/tamboto2000/otaqku-tasks/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("Error on loading config: %v", err))
		os.Exit(1)
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error on creating app: %v", err))
		os.Exit(1)
	}

	// Run database migration
	if err := app.RunDatabaseMigration(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Run HTTP server
	go func() {
		slog.Info("HTTP server started")
		err := app.RunHTTPServer()
		if err != nil {
			if err != http.ErrServerClosed {
				slog.Error(fmt.Sprintf("HTTP server stopped with error: %v", err))
			}
		}
	}()

	stopAppOnOsSignal(app)
}

func stopAppOnOsSignal(a *app.App) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	<-osSignal

	err := a.Shutdown()
	if err != nil {
		slog.Error(fmt.Sprintf("error on app shutdown: %v", err))
	}

	slog.Warn("HTTP server stopped")
}
