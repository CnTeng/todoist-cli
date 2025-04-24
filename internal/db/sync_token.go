package db

import (
	"context"
	"database/sql"
)

const (
	syncTokenStoreQuery = `
		INSERT INTO
			sync_token (id, token)
		VALUES
			(1, ?)
		ON CONFLICT (id) DO UPDATE
		SET
			token = excluded.token`

	syncTokenGetQuery = `SELECT token FROM sync_token WHERE id = 1`
)

func (db *DB) storeSyncToken(tx *sql.Tx, token string) error {
	_, err := tx.Exec(syncTokenStoreQuery, token)
	return err
}

func (db *DB) getSyncToken(ctx context.Context) (*string, error) {
	var token string
	if err := db.QueryRowContext(ctx, syncTokenGetQuery).Scan(&token); err != nil {
		return nil, err
	}
	return &token, nil
}
