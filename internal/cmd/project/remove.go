package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewRemoveCmd(f *util.Factory) *cli.Command {
	ids := []string{}
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
				Values: &ids,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			params := make([]*sync.ProjectDeleteArgs, 0, len(ids))
			for _, id := range ids {
				params = append(params, &sync.ProjectDeleteArgs{ID: id})
			}

			if _, err := f.Call(ctx, daemon.ProjectRemove, params); err != nil {
				return err
			}

			fmt.Printf("Project deleted: %s\n", strings.Join(ids, ", "))
			return nil
		},
	}
}
