package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync"
)

const (
	userStoreQuery = `
		INSERT INTO
			users (id, data)
		VALUES
			(?, ?)
		ON CONFLICT (id) DO UPDATE
		SET
			data = excluded.data`
	userDeleteQuery = `DELETE FROM users WHERE id = ?`

	userGetQuery = `SELECT data FROM users LIMIT 1`
)

func (db *DB) storeUser(ctx context.Context, tx *sql.Tx, user *sync.User) error {
	if user.IsDeleted {
		_, err := tx.ExecContext(ctx, userDeleteQuery, user.ID)
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, userStoreQuery, user.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUser(ctx context.Context) (*sync.User, error) {
	u := &sync.User{}
	return u, db.withTx(func(tx *sql.Tx) error {
		var err error
		u, err = getItem[sync.User](ctx, tx, userGetQuery)
		return err
	})
}
