package db

import (
	"context"
	"database/sql"
	_ "embed"
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

func NewDB(path string) (*DB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	dbPath := fmt.Sprintf("file:%s?_time_format=sqlite", path)
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
