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

	sectionListTemplate = `
		SELECT
			json_patch(
				section,
				json_object('project_name', project ->> 'name')
			) AS data
		FROM
			sections_view
		WHERE
			TRUE {{ . }}
		ORDER BY
			section ->> 'is_archived' ASC,
			project ->> 'child_order' ASC,
			section ->> 'section_order' ASC`
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
	sectionGetQuery, args, err := db.buildListQuery(
		sectionListTemplate,
		filters{"id": {Query: "id = ?", Arg: id}},
	)
	if err != nil {
		return nil, err
	}

	s := &model.Section{}
	return s, db.withTx(func(tx *sql.Tx) error {
		var err error
		s, err = getItem[model.Section](ctx, tx, sectionGetQuery, args...)
		return err
	})
}

func parseSectionFilters(args *model.SectionListArgs) filters {
	filters := filters{
		"is_archived": {Query: "section ->> 'is_archived' = false"},
	}

	if args == nil {
		return filters
	}

	if args.ProjectID != "" {
		filters["project.id"] = &filter{
			Query: "project ->> 'id' = ?",
			Arg:   args.ProjectID,
		}
	}

	if args.All {
		delete(filters, "is_archived")
	}

	if args.OnlyArchived {
		filters["is_archived"] = &filter{Query: "section ->> 'is_archived' = true"}
	}

	return filters
}

func (db *DB) ListSections(ctx context.Context, args *model.SectionListArgs) ([]*model.Section, error) {
	filters := parseSectionFilters(args)
	query, qargs, err := db.buildListQuery(sectionListTemplate, filters)
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
