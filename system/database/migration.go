package database

import (
	"database/sql"
	"github/anurag/altar-be/system/config"
	"github.com/pressly/goose/v3"
)

func RunMigrations(cfg *config.Config) error {
	db, err := sql.Open("pgx", cfg.PG_URL)
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
