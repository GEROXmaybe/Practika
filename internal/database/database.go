package database

import (
	"database/sql"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

func Init(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	migrationSQL, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return nil, err
	}

	slog.Info("база данных успешно инициализирована")
	return db, nil
}
