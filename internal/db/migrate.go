package db

import (
	"database/sql"
	_ "embed"
)

var (
	//go:embed migration/create_table.sql
	createTableQuery string

	//go:embed migration/create_view.sql
	createViewQuery string
)

var migrations = []string{
	createTableQuery,
	createViewQuery,
}

func (db *DB) Migrate() error {
	return db.withTx(func(tx *sql.Tx) error {
		for _, migration := range migrations {
			if _, err := tx.Exec(migration); err != nil {
				return err
			}
		}
		return nil
	})
}
