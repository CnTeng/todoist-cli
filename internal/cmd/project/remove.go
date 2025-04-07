package project

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewRemoveCmd(f *util.Factory) *cli.Command {
	params := []string{}
	return &cli.Command{
		Name:        "remove",
		Aliases:     []string{"rm"},
		Usage:       "Remove a project",
		Description: "Remove a project in todoist",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:   "id",
				Min:    1,
				Max:    -1,
				Values: &params,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			for _, id := range params {
				if _, err := f.Call(ctx, daemon.ProjectRemove, &sync.ProjectDeleteArgs{ID: id}); err != nil {
					return err
				}
				fmt.Printf("Project deleted: %s\n", id)
			}

			return nil
		},
	}
}
