package daemon

import (
	"context"
	"time"

	"github.com/CnTeng/todoist-api-go/rest"
	"github.com/CnTeng/todoist-api-go/todoist"
	"github.com/CnTeng/todoist-cli/internal/model"
)

func (d *Daemon) sync(ctx context.Context, args *model.SyncArgs) error {
	if _, err := d.client.SyncWithAutoToken(ctx, args.Force); err != nil {
		return err
	}

	if !args.All {
		return nil
	}

	taskSvc := todoist.NewTaskService(d.client)
	if _, err := taskSvc.GetCompletedTasksByCompletionDate(ctx, &rest.TaskGetCompletedByCompletionDateParams{
		Since: args.Since,
		Until: time.Now(),
	}); err != nil {
		return err
	}

	projectSvc := todoist.NewProjectService(d.client)
	if _, err := projectSvc.GetArchivedProjects(ctx, &rest.ProjectGetArchivedParams{}); err != nil {
		return err
	}

	return nil
}
