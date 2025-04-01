package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync"
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

	projectGetQuery  = `SELECT data FROM projects WHERE id = ?`
	projectListQuery = `
		SELECT
			data
		FROM
			projects
		ORDER BY
			data ->> 'inbox_project' DESC,
			data ->> 'child_order' ASC`
)

func (db *DB) StoreProject(project *sync.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	if _, err := db.Exec(projectStoreQuery, project.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetProject(id string) (*sync.Project, error) {
	var data []byte
	if err := db.QueryRow(projectGetQuery, id).Scan(&data); err != nil {
		return nil, err
	}

	project := &sync.Project{}
	if err := json.Unmarshal(data, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (db *DB) listProjects(ctx context.Context, tx *sql.Tx) ([]*sync.Project, error) {
	rows, err := tx.QueryContext(ctx, projectListQuery)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	ps := []*sync.Project{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		p := &sync.Project{}
		if err := json.Unmarshal(data, p); err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}

	return ps, nil
}

func (db *DB) ListProjects(ctx context.Context) ([]*sync.Project, error) {
	ps := []*sync.Project{}

	if err := db.withTx(func(tx *sql.Tx) error {
		var err error
		ps, err = db.listProjects(ctx, tx)
		return err
	}); err != nil {
		return nil, err
	}

	return ps, nil
}
