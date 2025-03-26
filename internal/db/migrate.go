package db

import (
	"database/sql"
	_ "embed"
)

//go:embed migrations/0001_create.sql
var m0001CreateQuery string

var migrations = []string{
	m0001CreateQuery,
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
