package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github/anurag/altar-be/system/config"
)

func RunMigrations(cfg *config.Config) error {
	db, err := sql.Open("pgx", cfg.Database.PostgresURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "system/migrations"); err != nil {
		return err
	}

	return nil
}
