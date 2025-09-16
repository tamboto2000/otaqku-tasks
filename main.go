package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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
	}
}
