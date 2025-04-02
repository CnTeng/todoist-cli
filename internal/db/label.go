package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync"
)

const (
	labelStoreQuery = `
		INSERT INTO
			labels (id, data)
		VALUES 
			(?, ?)
		ON CONFLICT (id) DO UPDATE
		SET
			data = excluded.data`
	labelDeleteQuery = `DELETE FROM labels WHERE id = ?`

	labelGetQuery  = `SELECT data FROM labels WHERE data ->> 'name' = ?`
	labelListQuery = `SELECT data FROM labels	ORDER BY data ->> 'item_order' ASC`
)

func (db *DB) storeLabel(tx *sql.Tx, label *sync.Label) error {
	if label.IsDeleted {
		_, err := tx.Exec(labelDeleteQuery, label.Name)
		return err
	}

	data, err := json.Marshal(label)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(labelStoreQuery, label.Name, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetLabel(ctx context.Context, name string) (*sync.Label, error) {
	l := &sync.Label{}
	return l, db.withTx(func(tx *sql.Tx) error {
		var err error
		l, err = getItem[sync.Label](ctx, tx, labelGetQuery, name)
		return err
	})
}

func (db *DB) ListLabels(ctx context.Context) ([]*sync.Label, error) {
	ls := []*sync.Label{}
	return ls, db.withTx(func(tx *sql.Tx) error {
		var err error
		ls, err = listItems[sync.Label](ctx, tx, labelListQuery)
		return err
	})
}
