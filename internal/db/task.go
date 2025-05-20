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

	taskListTemplate = `
		SELECT
			json_patch(
				task,
				json_object(
					'project_name',
					project ->> 'name',
					'project_color',
					project ->> 'color',
					'section_name',
					section ->> 'name'
				)
			) AS data
		FROM
			tasks_view
		WHERE
			TRUE {{ . }}
		ORDER BY
			project ->> 'inbox_project' DESC,
			project ->> 'child_order' ASC,
			task ->> 'checked' ASC,
			task ->> 'child_order' ASC`
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
	taskGetQuery, args, err := db.buildListQuery(
		taskListTemplate,
		filters{"id": {Query: "id = ?", Arg: id}},
	)
	if err != nil {
		return nil, err
	}

	t := &model.Task{}
	return t, db.withTx(func(tx *sql.Tx) error {
		var err error
		t, err = getItem[model.Task](ctx, tx, taskGetQuery, args...)
		return err
	})
}

func (db *DB) listTasks(ctx context.Context, tx *sql.Tx, filters filters) ([]*model.Task, error) {
	query, args, err := db.buildListQuery(taskListTemplate, filters)
	if err != nil {
		return nil, err
	}

	ts, err := listItems[model.Task](ctx, tx, query, args...)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		filters["task.parent_id"] = &filter{
			Query: "task ->> 'parent_id' = ?",
			Arg:   t.ID,
		}

		subTasks, err := db.listTasks(ctx, tx, filters)
		if err != nil {
			return nil, err
		}
		t.SubTasks = subTasks
	}

	return ts, nil
}

func (db *DB) ListTasks(ctx context.Context, args *model.TaskListArgs) ([]*model.Task, error) {
	filters := filters{
		"project.is_archived": {Query: "project ->> 'is_archived' = false"},
		"task.checked":        {Query: "task ->> 'checked' = false"},
		"task.parent_id":      {Query: "task ->> 'parent_id' IS NULL"},
	}
	if args != nil {
		if args.Completed {
			delete(filters, "task.checked")
		}

		if args.ProjectID != "" {
			filters["project.id"] = &filter{
				Query: "project ->> 'id' = ?",
				Arg:   args.ProjectID,
			}
		}
	}

	ts := []*model.Task{}
	return ts, db.withTx(func(tx *sql.Tx) error {
		var err error
		ts, err = db.listTasks(ctx, tx, filters)
		return err
	})
}
