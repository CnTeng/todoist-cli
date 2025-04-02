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

func (db *DB) getLabel(ctx context.Context, tx *sql.Tx, name string) (*sync.Label, error) {
	var data []byte
	if err := tx.QueryRowContext(ctx, labelGetQuery, name).Scan(&data); err != nil {
		return nil, err
	}

	label := &sync.Label{}
	if err := json.Unmarshal(data, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (db *DB) listGetLabel(ctx context.Context, tx *sql.Tx) ([]*sync.Label, error) {
	rows, err := tx.QueryContext(ctx, labelListQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ls := []*sync.Label{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		label := &sync.Label{}
		if err := json.Unmarshal(data, label); err != nil {
			return nil, err
		}

		ls = append(ls, label)
	}

	return ls, nil
}

func (db *DB) ListLabels(ctx context.Context) ([]*sync.Label, error) {
	ls := []*sync.Label{}

	if err := db.withTx(func(tx *sql.Tx) error {
		var err error
		ls, err = db.listGetLabel(ctx, tx)
		return err
	}); err != nil {
		return nil, err
	}

	return ls, nil
}
