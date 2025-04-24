package daemon

import (
	"context"
	"time"

	"github.com/CnTeng/todoist-api-go/rest"
	"github.com/CnTeng/todoist-api-go/ws"
)

func (d *Daemon) HandleMessage(ctx context.Context, msg ws.Message) error {
	switch msg {
	case ws.Connected:
		fallthrough
	case ws.SyncNeeded:
		d.log.Println("sync needed")
		if _, err := d.client.SyncWithAutoToken(ctx, false); err != nil {
			return err
		}
	}
	return nil
}

type SyncArgs struct {
	Force bool
	All   bool
	Since time.Time
}

func (d *Daemon) sync(ctx context.Context, args *SyncArgs) error {
	if _, err := d.client.SyncWithAutoToken(ctx, args.Force); err != nil {
		return err
	}

	if !args.All {
		return nil
	}

	ps, err := d.db.ListProjects(ctx)
	if err != nil {
		return err
	}

	for _, p := range ps {
		params := &rest.TaskGetCompletedByCompletionDateParams{
			Since:     args.Since,
			Until:     time.Now(),
			ProjectID: &p.ID,
		}

		resp, err := d.client.GetTasksCompletedByCompletionDate(ctx, params)
		if err != nil {
			return err
		}

		for resp.NextCursor != nil {
			params.Cursor = resp.NextCursor
			resp, err = d.client.GetTasksCompletedByCompletionDate(ctx, params)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
