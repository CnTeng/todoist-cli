package db

import (
	"context"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync/v9"
	"github.com/CnTeng/todoist-cli/internal/model"
)

const (
	storeTaskQuery = `INSERT INTO tasks (id, data) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET data = excluded.data`
	getTaskQuery   = `SELECT data FROM tasks WHERE id = ?`
	listTasksQuery = `SELECT data FROM tasks`
)

func (db *DB) StoreTask(task *sync.Item) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if _, err := db.Exec(storeTaskQuery, task.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetTask(ctx context.Context, id string) (*model.Task, error) {
	var data []byte
	if err := db.QueryRow(getTaskQuery, id).Scan(&data); err != nil {
		return nil, err
	}

	t := &model.Task{Item: &sync.Item{}}
	if err := json.Unmarshal(data, t); err != nil {
		return nil, err
	}

	project, err := db.GetProject(t.ProjectID)
	if err != nil {
		return nil, err
	}

	t.Project = project.Name

	return t, nil
}

func (db *DB) ListTasks(ctx context.Context) ([]*model.Task, error) {
	rows, err := db.Query(listTasksQuery)
	if err != nil {
		return nil, err
	}

	ts := []*model.Task{}

	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		t := &model.Task{Item: &sync.Item{}}
		if err := json.Unmarshal(data, t); err != nil {
			return nil, err
		}

		project, err := db.GetProject(t.ProjectID)
		if err != nil {
			return nil, err
		}
		t.Project = project.Name

		for _, ln := range t.Item.Labels {
			label, err := db.GetLabelByName(ln)
			if err != nil {
				return nil, err
			}
			t.Labels = append(t.Labels, label)
		}

		ts = append(ts, t)
	}

	return ts, nil
}
