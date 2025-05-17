package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"text/template"

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

	taskGetQuery          = `SELECT data FROM tasks_view WHERE id = ?`
	taskListQueryTemplate = `
		SELECT
			data
		FROM
			tasks_view
		WHERE
			data ->> '$.project.is_archived' = false
			{{ . }}
		ORDER BY
			data ->> '$.project.inbox_project' DESC,
			data ->> '$.project.child_order' ASC,
			data ->> 'checked' ASC,
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

func (db *DB) GetTask(ctx context.Context, id string) (*model.Task, error) {
	t := &model.Task{}
	return t, db.withTx(func(tx *sql.Tx) error {
		var err error
		t, err = getItem[model.Task](ctx, tx, taskGetQuery, id)
		return err
	})
}

type taskListCondition struct {
	Query string
	Args  []any
}

func (db *DB) buildTaskListQuery(conds map[string]*taskListCondition) (string, error) {
	t, err := template.New("taskListQuery").Parse(taskListQueryTemplate)
	if err != nil {
		return "", err
	}

	var cond string
	for _, c := range conds {
		cond += " AND " + c.Query
	}

	b := &strings.Builder{}
	if err := t.Execute(b, cond); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (db *DB) listTasks(ctx context.Context, tx *sql.Tx, conds map[string]*taskListCondition) ([]*model.Task, error) {
	taskListQuery, err := db.buildTaskListQuery(conds)
	if err != nil {
		return nil, err
	}

	args := []any{}
	for _, cond := range conds {
		args = append(args, cond.Args...)
	}

	ts, err := listItems[model.Task](ctx, tx, taskListQuery, args...)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		conds["parent_id"] = &taskListCondition{
			Query: "data ->> 'parent_id' = ?",
			Args:  []any{t.ID},
		}

		subTasks, err := db.listTasks(ctx, tx, conds)
		if err != nil {
			return nil, err
		}
		t.SubTasks = subTasks
	}

	return ts, nil
}

func (db *DB) ListTasks(ctx context.Context, all bool) ([]*model.Task, error) {
	ts := []*model.Task{}
	conds := map[string]*taskListCondition{
		"checked":   {Query: "data ->> 'checked' = false"},
		"parent_id": {Query: "data ->> 'parent_id' IS NULL"},
	}

	if all {
		delete(conds, "checked")
	}
	return ts, db.withTx(func(tx *sql.Tx) error {
		var err error
		ts, err = db.listTasks(ctx, tx, conds)
		return err
	})
}
