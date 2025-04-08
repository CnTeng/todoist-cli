package task

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/urfave/cli/v3"
)

func NewListCmd(f *util.Factory) *cli.Command {
	params := &view.TaskViewConfig{}
	return &cli.Command{
		Name:                   "list",
		Aliases:                []string{"ls"},
		Usage:                  "List tasks",
		Description:            "List tasks in todoist",
		Category:               "task",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				Usage:       "list all tasks include completed",
				HideDefault: true,
				OnlyOnce:    true,
				Destination: &params.Completed,
			},
			&cli.BoolFlag{
				Name:        "description",
				Aliases:     []string{"d"},
				Usage:       "list tasks include description",
				HideDefault: true,
				OnlyOnce:    true,
				Destination: &params.Description,
			},
			&cli.BoolFlag{
				Name:        "tree",
				Aliases:     []string{"t"},
				Usage:       "list tasks in tree format",
				HideDefault: true,
				OnlyOnce:    true,
				Destination: &params.Tree,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			result := []*model.Task{}
			if err := f.CallResult(ctx, daemon.TaskList, params, &result); err != nil {
				return err
			}

			v := view.NewTaskView(result, f.IconConfig.Icons, params)
			fmt.Print(v.Render())
			return nil
		},
	}
}
