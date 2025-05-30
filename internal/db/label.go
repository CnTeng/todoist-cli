package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/model"
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
	labelDeleteQuery = `DELETE FROM labels WHERE data ->> 'id' = ?`

	labelGetQuery  = `SELECT data FROM labels_view WHERE id = ?`
	labelListQuery = `
		SELECT
			data
		FROM
			labels_view
		ORDER BY
			data ->> 'is_favorite' DESC,
			data ->> 'is_shared' ASC,
			data ->> 'item_order' ASC`
)

func (db *DB) storeLabel(ctx context.Context, tx *sql.Tx, label *sync.Label) error {
	if label.IsDeleted {
		_, err := tx.ExecContext(ctx, labelDeleteQuery, label.ID)
		return err
	}

	data, err := json.Marshal(&model.Label{Label: label})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, labelStoreQuery, label.Name, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetLabel(ctx context.Context, name string) (*model.Label, error) {
	l := &model.Label{}
	return l, db.withTx(func(tx *sql.Tx) error {
		var err error
		l, err = getItem[model.Label](ctx, tx, labelGetQuery, name)
		return fmt.Errorf("failed to get label %s: %w", name, err)
	})
}

func (db *DB) ListLabels(ctx context.Context) ([]*model.Label, error) {
	ls := []*model.Label{}
	return ls, db.withTx(func(tx *sql.Tx) error {
		var err error
		ls, err = listItems[model.Label](ctx, tx, labelListQuery)
		return fmt.Errorf("failed to list labels: %w", err)
	})
}
