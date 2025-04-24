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

	taskListAllSubQuery = `
		SELECT
			data
		FROM
			tasks
		WHERE
			data ->> 'parent_id' = ?
		ORDER BY
			(data ->> 'completed_at' IS NOT NULL) ASC,
			data ->> 'child_order' ASC`

	taskListUndoneSubQuery = `
		SELECT
			data
		FROM
			tasks
		WHERE
			data ->> 'parent_id' = ?
			AND data ->> 'completed_at' IS NULL
		ORDER BY
			data ->> 'child_order' ASC`

	taskListAllByProjectQuery = `
		SELECT
			data
		FROM
			tasks
		WHERE
			data ->> 'parent_id' IS NULL
			AND data ->> 'project_id' = ?
		ORDER BY
			(data ->> 'completed_at' IS NOT NULL) ASC,
			data ->> 'child_order' ASC`

	taskListUndoneByProjectQuery = `
		SELECT
			data
		FROM
			tasks
		WHERE
			data ->> 'parent_id' IS NULL
			AND data ->> 'completed_at' IS NULL
			AND data ->> 'project_id' = ?
		ORDER BY
			data ->> 'child_order' ASC`
)

func (db *DB) storeTask(ctx context.Context, tx *sql.Tx, task *sync.Task) error {
	if task.IsDeleted {
		_, err := tx.ExecContext(ctx, taskDeleteQuery, task.ID)
		return err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, taskStoreQuery, task.ID, data); err != nil {
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

	t := &model.Task{Task: &sync.Task{}}
	if err := json.Unmarshal(tdata, t); err != nil {
		return nil, err
	}

	p := &sync.Project{}
	if err := json.Unmarshal(pdata, p); err != nil {
		return nil, err
	}
	t.Project = p

	for _, ln := range t.Task.Labels {
		label, err := getItem[sync.Label](ctx, tx, labelGetQuery, ln)
		if err != nil {
			return nil, err
		}
		t.Labels = append(t.Labels, label)
	}

	return t, nil
}

func (db *DB) listSubTasks(ctx context.Context, tx *sql.Tx, query string, task *model.Task) error {
	rows, err := tx.QueryContext(ctx, query, task.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	task.SubTasks = []*model.Task{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}

		st := &model.Task{Task: &sync.Task{}, Project: task.Project}
		if err := json.Unmarshal(data, st.Task); err != nil {
			return err
		}

		if err := db.listSubTasks(ctx, tx, query, st); err != nil {
			return err
		}

		for _, ln := range st.Task.Labels {
			label, err := getItem[sync.Label](ctx, tx, labelGetQuery, ln)
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

func (db *DB) listTasksByProject(ctx context.Context, tx *sql.Tx, project *sync.Project, all bool) ([]*model.Task, error) {
	taskListByProjectQuery := taskListUndoneByProjectQuery
	taskListSubQuery := taskListUndoneSubQuery
	if all {
		taskListByProjectQuery = taskListAllByProjectQuery
		taskListSubQuery = taskListAllSubQuery
	}

	rows, err := tx.QueryContext(ctx, taskListByProjectQuery, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ts := []*model.Task{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		t := &model.Task{Task: &sync.Task{}, Project: project}
		if err := json.Unmarshal(data, t.Task); err != nil {
			return nil, err
		}

		if err := db.listSubTasks(ctx, tx, taskListSubQuery, t); err != nil {
			return nil, err
		}

		for _, ln := range t.Task.Labels {
			label, err := getItem[sync.Label](ctx, tx, labelGetQuery, ln)
			if err != nil {
				// TODO: shared labels
				label = &sync.Label{Name: ln, Color: sync.Grey}
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

func (db *DB) ListTasks(ctx context.Context, all bool) ([]*model.Task, error) {
	ts := []*model.Task{}

	if err := db.withTx(func(tx *sql.Tx) error {
		ps, err := listItems[sync.Project](ctx, tx, projectListQuery)
		if err != nil {
			return err
		}

		for _, p := range ps {
			t, err := db.listTasksByProject(ctx, tx, p, all)
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
