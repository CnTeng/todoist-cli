package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync"
)

const (
	filterStoreQuery = `
		INSERT INTO
			filters (id, data)
		VALUES
			(?, ?)
		ON CONFLICT (id) DO UPDATE
		SET
			data = excluded.data`
	filterDeleteQuery = `DELETE FROM filters WHERE id = ?`

	filterGetQuery  = `SELECT data FROM filters WHERE id = ?`
	filterListQuery = `
		SELECT
			data
		FROM
			filters
		ORDER BY
			data ->> 'item_order' ASC`
)

func (db *DB) storeFilter(ctx context.Context, tx *sql.Tx, filter *sync.Filter) error {
	if filter.IsDeleted {
		_, err := tx.ExecContext(ctx, filterDeleteQuery, filter.ID)
		return err
	}

	data, err := json.Marshal(filter)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, filterStoreQuery, filter.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetFilter(ctx context.Context, id string) (*sync.Filter, error) {
	f := &sync.Filter{}
	return f, db.withTx(func(tx *sql.Tx) error {
		var err error
		f, err = getItem[sync.Filter](ctx, tx, filterGetQuery, id)
		return err
	})
}

func (db *DB) ListFilters(ctx context.Context) ([]*sync.Filter, error) {
	fs := []*sync.Filter{}
	return fs, db.withTx(func(tx *sql.Tx) error {
		var err error
		fs, err = listItems[sync.Filter](ctx, tx, filterListQuery)
		return err
	})
}
