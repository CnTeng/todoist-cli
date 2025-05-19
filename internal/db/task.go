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

	taskGetQuery          = `SELECT data FROM tasks_view WHERE id = ?`
	taskListQueryTemplate = `
		SELECT
			data
		FROM
			tasks_view
		WHERE
			TRUE {{ . }}
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

func (db *DB) listTasks(ctx context.Context, tx *sql.Tx, conds listConditions) ([]*model.Task, error) {
	query, args, err := db.buildListQuery(taskListQueryTemplate, conds)
	if err != nil {
		return nil, err
	}

	ts, err := listItems[model.Task](ctx, tx, query, args...)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		conds["parent_id"] = &listCondition{
			Query: "data ->> 'parent_id' = ?",
			Arg:   t.ID,
		}

		subTasks, err := db.listTasks(ctx, tx, conds)
		if err != nil {
			return nil, err
		}
		t.SubTasks = subTasks
	}

	return ts, nil
}

func (db *DB) ListTasks(ctx context.Context, args *model.TaskListArgs) ([]*model.Task, error) {
	ts := []*model.Task{}
	conds := listConditions{
		"project.is_archived": {Query: "data ->> '$.project.is_archived' = false"},
		"checked":             {Query: "data ->> 'checked' = false"},
		"parent_id":           {Query: "data ->> 'parent_id' IS NULL"},
	}
	if args != nil && args.Completed {
		delete(conds, "checked")
	}

	return ts, db.withTx(func(tx *sql.Tx) error {
		var err error
		ts, err = db.listTasks(ctx, tx, conds)
		return err
	})
}
