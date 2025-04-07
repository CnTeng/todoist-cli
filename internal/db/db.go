package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const (
	cleanQuery = `
		PRAGMA foreign_keys = OFF;
		DELETE FROM tasks;
		DELETE FROM projects;
		DELETE FROM labels;
		PRAGMA foreign_keys = ON;
	`
)

type DB struct {
	*sql.DB
}

type Config struct {
	Path string `toml:"path"`
}

func NewDB(config *Config) (*DB, error) {
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		file, err := os.Create(config.Path)
		if err != nil {
			return nil, err
		}

		if err := file.Close(); err != nil {
			return nil, err
		}
	}

	dbPath := fmt.Sprintf("file:%s?_time_format=sqlite", config.Path)
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err = database.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	if _, err = database.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, err
	}

	if _, err = database.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (db *DB) clean(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, cleanQuery); err != nil {
		return err
	}
	return nil
}

func (db *DB) withTx(fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
