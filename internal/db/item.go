package db

import (
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync/v9"
	"github.com/CnTeng/todoist-cli/internal/model"
)

const (
	storeItemQuery = `INSERT INTO items (id, data) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET data = excluded.data`
	getItemQuery   = `SELECT data FROM items WHERE id = ?`
	listItemsQuery = `SELECT data FROM items`
)

func (db *DB) StoreItem(item *sync.Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	if _, err := db.Exec(storeItemQuery, item.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetTask(id string) (*model.Item, error) {
	var data []byte
	if err := db.QueryRow(getItemQuery, id).Scan(&data); err != nil {
		return nil, err
	}

	task := &model.Item{Item: &sync.Item{}}
	if err := json.Unmarshal(data, task); err != nil {
		return nil, err
	}

	project, err := db.GetProject(task.ProjectID)
	if err != nil {
		return nil, err
	}

	task.Project = project.Name

	return task, nil
}

func (db *DB) ListTasks() ([]*model.Item, error) {
	rows, err := db.Query(listItemsQuery)
	if err != nil {
		return nil, err
	}

	ts := []*model.Item{}

	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		t := &model.Item{Item: &sync.Item{}}
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
