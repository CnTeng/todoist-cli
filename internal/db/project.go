package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/model"
)

const (
	projectStoreQuery = `
		INSERT INTO
			projects (id, data)
		VALUES
			(?, ?)
		ON CONFLICT (id) DO UPDATE
		SET
			data = excluded.data`
	projectDeleteQuery = `DELETE FROM projects WHERE id = ?`

	projectGetQuery     = `SELECT data FROM projects WHERE id = ?`
	projectListTemplate = `
		SELECT
			data
		FROM
			projects
		WHERE
			TRUE {{ . }}
		ORDER BY
			data ->> 'inbox_project' DESC,
			data ->> 'is_favorite' DESC,
			data ->> 'is_archived' ASC,
			data ->> 'child_order' ASC`
)

func (db *DB) storeProject(ctx context.Context, tx *sql.Tx, project *sync.Project) error {
	if project.IsDeleted {
		_, err := tx.ExecContext(ctx, projectDeleteQuery, project.ID)
		return err
	}

	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, projectStoreQuery, project.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetProject(ctx context.Context, id string) (*sync.Project, error) {
	p := &sync.Project{}
	return p, db.withTx(func(tx *sql.Tx) error {
		var err error
		p, err = getItem[sync.Project](ctx, tx, projectGetQuery, id)
		return err
	})
}

func (db *DB) ListProjects(ctx context.Context, args *model.ProjectListArgs) ([]*sync.Project, error) {
	filters := filters{
		"is_archived": {Query: "data ->> 'is_archived' = false"},
	}
	if args != nil && args.Archived {
		delete(filters, "is_archived")
	}

	query, qargs, err := db.buildListQuery(projectListTemplate, filters)
	if err != nil {
		return nil, err
	}

	ps := []*sync.Project{}
	return ps, db.withTx(func(tx *sql.Tx) error {
		var err error
		ps, err = listItems[sync.Project](ctx, tx, query, qargs...)
		return err
	})
}
