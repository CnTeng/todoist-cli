package daemon

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-api-go/todoist"
	"github.com/CnTeng/todoist-cli/internal/model"
)

func (d *Daemon) updateLabel(ctx context.Context, args *model.LabelUpdateArgs) error {
	label, err := d.db.GetLabel(ctx, args.Name)
	if err != nil {
		return fmt.Errorf("label %s not found: %w", args.Name, err)
	}

	svc := todoist.NewLabelService(d.client)

	if label.IsShared {
		if args.Args.Name == nil {
			return fmt.Errorf("label %s is shared, name is required", label.Name)
		}

		_, err := svc.RenameSharedLabel(ctx, &sync.LabelRenameSharedArgs{
			NameOld: label.Name,
			NameNew: *args.Args.Name,
		})
		return err
	} else {
		_, err := svc.UpdateLabel(ctx, &args.Args)
		return err
	}
}

func (d *Daemon) deleteLabels(ctx context.Context, args []*model.LabelDeleteArgs) error {
	cmds := make(sync.Commands, len(args))

	for i, arg := range args {
		label, err := d.db.GetLabel(ctx, arg.Name)
		if err != nil {
			return fmt.Errorf("label %s not found: %w", arg.Name, err)
		}

		var cmd *sync.Command
		if label.IsShared {
			cmd = sync.NewCommand(&sync.LabelDeleteSharedArgs{
				Name: label.Name,
			})
		} else {
			cmd = sync.NewCommand(&sync.LabelDeleteArgs{
				ID: label.ID,
			})
		}
		cmds[i] = cmd
	}

	_, err := d.client.ExecuteCommands(ctx, cmds)
	return err
}
