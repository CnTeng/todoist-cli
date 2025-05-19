package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"text/template"

	_ "modernc.org/sqlite"
)

const (
	cleanQuery = `
		PRAGMA foreign_keys = OFF;
		DELETE FROM tasks;
		DELETE FROM projects;
		DELETE FROM labels;
		DELETE FROM users;
		PRAGMA foreign_keys = ON;
	`
)

type DB struct {
	*sql.DB
}

func NewDB(path string) (*DB, error) {
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

type listCondition struct {
	Query string
	Arg   any
}

type listConditions map[string]*listCondition

func (db *DB) buildListQuery(query string, conds listConditions) (string, []any, error) {
	t, err := template.New("listQuery").Parse(query)
	if err != nil {
		return "", nil, err
	}

	cb := &strings.Builder{}
	args := []any{}
	for _, c := range conds {
		cb.WriteString(" AND ")
		cb.WriteString(c.Query)

		if c.Arg != nil {
			args = append(args, c.Arg)
		}
	}

	b := &strings.Builder{}
	if err := t.Execute(b, cb.String()); err != nil {
		return "", nil, err
	}

	return b.String(), args, nil
}
