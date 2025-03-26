package db

import (
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync/v9"
)

const (
	storeProjectQuery = `INSERT INTO projects (id, data) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET data = excluded.data`
	getProjectQuery   = `SELECT data FROM projects WHERE id = ?`
	listProjectsQuery = `SELECT data FROM projects`
)

func (db *DB) StoreProject(project *sync.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	if _, err := db.Exec(storeProjectQuery, project.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetProject(id string) (*sync.Project, error) {
	var data []byte
	if err := db.QueryRow(getProjectQuery, id).Scan(&data); err != nil {
		return nil, err
	}

	project := &sync.Project{}
	if err := json.Unmarshal(data, project); err != nil {
		return nil, err
	}

	return project, nil
}
