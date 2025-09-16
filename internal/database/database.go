package database

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/tamboto2000/otaqku-tasks/internal/config"
	"github.com/vinovest/sqlx"
)

func Connect(cfg config.Database) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode,
	)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error on opening connection to database: %v", err)
	}

	return db, nil
}

func RunMigration(ctx context.Context, cfg config.Database, db *sqlx.DB) error {
	err := goose.UpContext(ctx, db.DB, cfg.MigrationDir)
	if err != nil {
		return fmt.Errorf("error on running database migration: %v", err)
	}

	return nil
}
