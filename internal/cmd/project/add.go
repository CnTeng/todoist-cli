package project

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewAddCmd(f *util.Factory) *cli.Command {
	params := &sync.ProjectAddArgs{}
	return &cli.Command{
		Name:        "add",
		Aliases:     []string{"a"},
		Usage:       "Add a project",
		Description: "Add a project to todoist",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "name",
				Destination: &params.Name,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Flags: []cli.Flag{
			newColorFlag(&params.Color),
			newFavoriteFlag(&params.IsFavorite),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.Call(ctx, daemon.ProjectAdd, params); err != nil {
				return err
			}

			fmt.Printf("Project added: %s\n", params.Name)

			return nil
		},
	}
}
