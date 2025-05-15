package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/model"
)

const (
	sectionStoreQuery = `
		INSERT INTO
			sections (id, data)
		VALUES
			(?, ?)
		ON CONFLICT (id) DO UPDATE
		SET
			data = excluded.data`
	sectionDeleteQuery = `DELETE FROM sections WHERE id = ?`

	sectionGetQuery  = `SELECT data FROM sections_view WHERE id = ?`
	sectionListQuery = `
		SELECT
			data
		FROM
			sections_view
		ORDER BY
			data ->> 'is_archived' ASC,
			data ->> 'project_name' ASC,
			data ->> 'section_order' ASC`
)

func (db *DB) storeSection(ctx context.Context, tx *sql.Tx, section *sync.Section) error {
	if section.IsDeleted {
		_, err := tx.ExecContext(ctx, sectionDeleteQuery, section.ID)
		return err
	}

	data, err := json.Marshal(section)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, sectionStoreQuery, section.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetSection(ctx context.Context, id string) (*model.Section, error) {
	s := &model.Section{}
	return s, db.withTx(func(tx *sql.Tx) error {
		var err error
		s, err = getItem[model.Section](ctx, tx, sectionGetQuery, id)
		return err
	})
}

func (db *DB) ListSections(ctx context.Context) ([]*model.Section, error) {
	ss := []*model.Section{}
	return ss, db.withTx(func(tx *sql.Tx) error {
		var err error
		ss, err = listItems[model.Section](ctx, tx, sectionListQuery)
		return err
	})
}
