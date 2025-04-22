package task

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewMoveCmd(f *util.Factory) *cli.Command {
	params := &sync.ItemMoveArgs{}
	return &cli.Command{
		Name:        "move",
		Aliases:     []string{"mv"},
		Usage:       "Move a task",
		Description: "Move a task in todoist",
		Category:    "task",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "id",
				Destination: &params.ID,
			},
		},
		Flags: []cli.Flag{
			newSectionFlag(&params.SectionID),
			newParentFlag(&params.ParentID),
			newProjectFlag(&params.ProjectID),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.Call(ctx, daemon.TaskMove, params); err != nil {
				return err
			}

			fmt.Printf("Task moved: %s\n", params.ID)
			return nil
		},
	}
}
