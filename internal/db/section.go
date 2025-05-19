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

	sectionGetQuery          = `SELECT data FROM sections_view WHERE id = ?`
	sectionListQueryTemplate = `
	SELECT
		json_patch(
			data,
			json_object('project_name', project ->> 'name')
		) AS data
	FROM
		sections_view
	WHERE
		TRUE {{ . }}
	ORDER BY
		data ->> 'is_archived' ASC,
		project ->> 'child_order' ASC,
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

func (db *DB) ListSections(ctx context.Context, args *model.SectionListArgs) ([]*model.Section, error) {
	conds := listConditions{}
	if args != nil && args.ProjectID != "" {
		conds["project_id"] = &listCondition{
			Query: "project ->> 'id' = ?",
			Arg:   args.ProjectID,
		}
	}

	query, qargs, err := db.buildListQuery(sectionListQueryTemplate, conds)
	if err != nil {
		return nil, err
	}

	ss := []*model.Section{}
	return ss, db.withTx(func(tx *sql.Tx) error {
		var err error
		ss, err = listItems[model.Section](ctx, tx, query, qargs...)
		return err
	})
}
