package task

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewReopenCmd(f *util.Factory) *cli.Command {
	params := []string{}
	return &cli.Command{
		Name:        "reopen",
		Aliases:     []string{"r"},
		Usage:       "Reopen a task",
		Description: "Reopen a task in todoist",
		Category:    "task",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "id",
				Min:         1,
				Max:         -1,
				Destination: &params,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			for _, id := range params {
				if _, err := f.Call(ctx, daemon.TaskReopen, &sync.TaskUncompleteArgs{ID: id}); err != nil {
					return err
				}
				fmt.Printf("Task reopen: %s\n", id)
			}

			return nil
		},
	}
}
