package project

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewModifyCmd(f *util.Factory) *cli.Command {
	params := &sync.ProjectUpdateArgs{}
	return &cli.Command{
		Name:        "modify",
		Aliases:     []string{"m"},
		Usage:       "Modify a task",
		Description: "Modify a task in todoist",
		Category:    "task",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "ID",
				Destination: &params.ID,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Flags: []cli.Flag{
			newNameFlag(&params.Name),
			newColorFlag(&params.Color),
			newFavoriteFlag(&params.IsFavorite),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.Call(ctx, daemon.ProjectModify, params); err != nil {
				return err
			}

			fmt.Printf("Project modified: %s\n", params.ID)
			return nil
		},
	}
}
