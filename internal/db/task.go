package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/model"
)

const (
	taskStoreQuery = `
		INSERT INTO
			tasks (id, data)
		VALUES
			(?, ?)
		ON CONFLICT (id) DO UPDATE
		SET
			data = excluded.data`
	taskDeleteQuery = `DELETE FROM tasks WHERE id = ?`

	taskGetQuery = `
		SELECT
			t.data AS task_data,
			p.data AS project_data
		FROM
			tasks t
			JOIN projects p ON p.id = t.data ->> 'project_id'
		WHERE
			id = ?`

	taskListSubQuery = `
		SELECT
			data
		FROM
			tasks
		WHERE
			data ->> 'parent_id' = ?
		ORDER BY
			data ->> 'child_order' ASC`

	taskListByProjectQuery = `
		SELECT
			data
		FROM
			tasks
		WHERE
			data ->> 'parent_id' IS NULL
			AND data ->> 'project_id' = ?
		ORDER BY
			data ->> 'child_order' ASC`
)

func (db *DB) storeTask(tx *sql.Tx, task *sync.Item) error {
	if task.IsDeleted {
		_, err := tx.Exec(taskDeleteQuery, task.ID)
		return err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(taskStoreQuery, task.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) getTask(ctx context.Context, tx *sql.Tx, id string) (*model.Task, error) {
	var tdata []byte
	var pdata []byte
	if err := db.QueryRowContext(ctx, taskGetQuery, id).Scan(&tdata, &pdata); err != nil {
		return nil, err
	}

	t := &model.Task{Item: &sync.Item{}}
	if err := json.Unmarshal(tdata, t); err != nil {
		return nil, err
	}

	p := &sync.Project{}
	if err := json.Unmarshal(pdata, p); err != nil {
		return nil, err
	}
	t.Project = p

	for _, ln := range t.Item.Labels {
		label, err := db.getLabelByName(ctx, tx, ln)
		if err != nil {
			return nil, err
		}
		t.Labels = append(t.Labels, label)
	}

	return t, nil
}

func (db *DB) listSubTasks(ctx context.Context, tx *sql.Tx, task *model.Task) error {
	rows, err := tx.QueryContext(ctx, taskListSubQuery, task.ID)
	if err != nil {
		return err
	}

	task.SubTasks = []*model.Task{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}

		st := &model.Task{Item: &sync.Item{}, Project: task.Project}
		if err := json.Unmarshal(data, st.Item); err != nil {
			return err
		}

		if err := db.listSubTasks(ctx, tx, st); err != nil {
			return err
		}

		for _, ln := range st.Item.Labels {
			label, err := db.getLabelByName(ctx, tx, ln)
			if err != nil {
				return err
			}
			st.Labels = append(st.Labels, label)
		}

		task.SubTaskStatus.Total++
		if st.CompletedAt != nil {
			task.SubTaskStatus.Completed++
		}

		task.SubTasks = append(task.SubTasks, st)
	}

	return nil
}

func (db *DB) listTasksByProject(ctx context.Context, tx *sql.Tx, project *sync.Project) ([]*model.Task, error) {
	rows, err := tx.QueryContext(ctx, taskListByProjectQuery, project.ID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	ts := []*model.Task{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		t := &model.Task{Item: &sync.Item{}, Project: project}
		if err := json.Unmarshal(data, t.Item); err != nil {
			return nil, err
		}

		if err := db.listSubTasks(ctx, tx, t); err != nil {
			return nil, err
		}

		for _, ln := range t.Item.Labels {
			label, err := db.getLabelByName(ctx, tx, ln)
			if err != nil {
				return nil, err
			}
			t.Labels = append(t.Labels, label)
		}

		ts = append(ts, t)
	}

	return ts, nil
}

func (db *DB) GetTask(ctx context.Context, id string) (*model.Task, error) {
	var t *model.Task
	return t, db.withTx(func(tx *sql.Tx) error {
		var err error
		t, err = db.getTask(ctx, tx, id)
		return err
	})
}

func (db *DB) ListTasks(ctx context.Context) ([]*model.Task, error) {
	ts := []*model.Task{}

	if err := db.withTx(func(tx *sql.Tx) error {
		ps, err := db.listProjects(ctx, tx)
		if err != nil {
			return err
		}

		for _, p := range ps {
			t, err := db.listTasksByProject(ctx, tx, p)
			if err != nil {
				return err
			}
			ts = append(ts, t...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ts, nil
}
